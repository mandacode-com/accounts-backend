package oauthauth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/auth/ent/authaccount"

	"mandacode.com/accounts/auth/internal/infra/oauthapi"
	dbmodels "mandacode.com/accounts/auth/internal/models/database"
	oauthmodels "mandacode.com/accounts/auth/internal/models/oauth"
	coderepo "mandacode.com/accounts/auth/internal/repository/code"
	dbrepo "mandacode.com/accounts/auth/internal/repository/database"
	tokenrepo "mandacode.com/accounts/auth/internal/repository/token"
	userrepo "mandacode.com/accounts/auth/internal/repository/user"
	oauthdto "mandacode.com/accounts/auth/internal/usecase/oauthauth/dto"
)

type LoginUsecase struct {
	authAccount      *dbrepo.AuthAccountRepository
	userService      *userrepo.UserServiceRepository
	token            *tokenrepo.TokenRepository
	loginCodeManager *coderepo.CodeManager
	oauthApiMap      map[authaccount.Provider]oauthapi.OAuthAPI
}

// createOAuth creates a new OAuth account in the database.
func (l *LoginUsecase) createOAuth(ctx context.Context, provider authaccount.Provider, userInfo *oauthmodels.UserInfo) (*dbmodels.SecureOAuthAuthAccount, error) {
	userID := uuid.New()
	initUser, err := l.userService.InitUser(ctx, userID)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to initialize user", errcode.ErrInternalFailure)
	}
	if initUser.UserId != userID.String() {
		l.userService.DeleteUser(ctx, userID)
		return nil, errors.New("user initialization failed", "User Initialization Error", errcode.ErrInternalFailure)
	}

	account, err := l.authAccount.CreateOAuthAuthAccount(
		ctx,
		&dbmodels.CreateOAuthAuthAccountInput{
			UserID:     userID,
			Provider:   provider,
			ProviderID: userInfo.ProviderID,
			Email:      userInfo.Email,
			IsVerified: userInfo.EmailVerified,
		},
	)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to create OAuth account", errcode.ErrInternalFailure)
	}
	return account, nil
}

// getAccessToken retrieves the access token from the OAuth API.
func (l *LoginUsecase) getAccessToken(ctx context.Context, provider authaccount.Provider, code string) (string, error) {
	api, ok := l.oauthApiMap[provider]
	if !ok {
		return "", errors.New(fmt.Sprintf("unsupported provider: %s", provider), "UnsupportedProvider", errcode.ErrInvalidInput)
	}
	accessToken, err := api.GetAccessToken(code)
	if err != nil {
		return "", errors.Upgrade(err, "Failed to get access token from OAuth provider", errcode.ErrUnauthorized)
	}
	return accessToken, nil
}

func (l *LoginUsecase) getOrCreateVerifiedUser(ctx context.Context, input oauthdto.LoginInput) (uuid.UUID, error) {
	var oauthAccessToken string
	if input.AccessToken == "" && input.Code != "" {
		var err error
		oauthAccessToken, err = l.getAccessToken(ctx, input.Provider, input.Code)
		if err != nil {
			return uuid.Nil, errors.Upgrade(err, "Failed to get access token", errcode.ErrUnauthorized)
		}
	} else if input.AccessToken != "" {
		oauthAccessToken = input.AccessToken
	} else {
		return uuid.Nil, errors.New("either access token or code must be provided", "Invalid Input", errcode.ErrInvalidInput)
	}

	userInfo, err := l.oauthApiMap[input.Provider].GetUserInfo(oauthAccessToken)
	if err != nil {
		return uuid.Nil, errors.Upgrade(err, "Failed to get user info from OAuth provider", errcode.ErrUnauthorized)
	}
	if userInfo == nil {
		return uuid.Nil, errors.New("user info is nil", "InvalidUserInfo", errcode.ErrInvalidInput)
	}

	var verified bool
	var userID uuid.UUID
	oauth, err := l.authAccount.GetOAuthAccountByProviderAndProviderID(ctx, input.Provider, userInfo.ProviderID)
	if err != nil {
		if errors.Is(err, errcode.ErrNotFound) {
			// User not found, create a new OAuth account
			newAccount, err := l.createOAuth(ctx, input.Provider, userInfo)
			if err != nil {
				return uuid.Nil, errors.Upgrade(err, "Failed to create OAuth account", errcode.ErrInternalFailure)
			}
			userID = newAccount.UserID
			verified = newAccount.IsVerified
		}
		return uuid.Nil, errors.Upgrade(err, "Failed to get OAuth account", errcode.ErrInternalFailure)
	} else {
		userID = oauth.UserID
		verified = oauth.IsVerified
	}

	if !verified {
		return uuid.Nil, errors.New("user is not verified", "Unauthorized", errcode.ErrUnauthorized)
	}

	return userID, nil
}

// GetLoginURL implements oauthdomain.LoginUsecase.
func (l *LoginUsecase) GetLoginURL(ctx context.Context, provider string) (loginURL string, err error) {
	api, ok := l.oauthApiMap[authaccount.Provider(provider)]
	if !ok {
		return "", errors.New("unsupported provider: "+provider, "Unsupported Provider", errcode.ErrInvalidInput)
	}
	loginURL = api.GetLoginURL()
	return loginURL, nil
}

// IssueLoginCode implements oauthdomain.LoginUsecase.
func (l *LoginUsecase) IssueLoginCode(ctx context.Context, input oauthdto.LoginInput) (code string, userID uuid.UUID, err error) {
	// Get or create verified user
	userID, err = l.getOrCreateVerifiedUser(ctx, input)
	if err != nil {
		return "", uuid.Nil, errors.Upgrade(err, "Failed to get or create verified user", errcode.ErrUnauthorized)
	}

	// Generate and store login code
	code, err = l.loginCodeManager.IssueCode(ctx, userID)
	if err != nil {
		return "", uuid.Nil, errors.Upgrade(err, "Failed to issue login code", errcode.ErrInternalFailure)
	}

	return code, userID, nil
}

// Login implements oauthdomain.LoginUsecase.
func (l *LoginUsecase) Login(ctx context.Context, input oauthdto.LoginInput) (accessToken string, refreshToken string, err error) {
	// Get or create verified user
	userID, err := l.getOrCreateVerifiedUser(ctx, input)
	if err != nil {
		return "", "", errors.Upgrade(err, "Failed to get or create verified user", errcode.ErrUnauthorized)
	}

	// Generate access and refresh tokens
	return l.issueToken(ctx, userID)
}

// VerifyLoginCode implements oauthdomain.LoginUsecase.
func (l *LoginUsecase) VerifyLoginCode(ctx context.Context, userID uuid.UUID, code string) (accessToken string, refreshToken string, err error) {
	// Validate code
	valid, err := l.loginCodeManager.ValidateCode(ctx, userID, code)
	if err != nil {
		return "", "", errors.Upgrade(err, "Failed to validate login code", errcode.ErrInternalFailure)
	}
	if !valid {
		return "", "", errors.New("login code is invalid or expired", "Failed to validate login code", errcode.ErrUnauthorized)
	}

	// Generate access and refresh tokens
	return l.issueToken(ctx, userID)
}

// issueToken generates access and refresh tokens for the user.
func (l *LoginUsecase) issueToken(ctx context.Context, userID uuid.UUID) (accessToken string, refreshToken string, err error) {
	// Generate access token
	accessToken, _, err = l.token.GenerateAccessToken(ctx, userID)
	if err != nil {
		return "", "", errors.Upgrade(err, "Failed to generate access token", errcode.ErrInternalFailure)
	}

	// Generate refresh token
	refreshToken, _, err = l.token.GenerateRefreshToken(ctx, userID)
	if err != nil {
		return "", "", errors.Upgrade(err, "Failed to generate refresh token", errcode.ErrInternalFailure)
	}

	return accessToken, refreshToken, nil
}

// NewLoginUsecase creates a new instance of LoginUsecase.
func NewLoginUsecase(
	authAccount *dbrepo.AuthAccountRepository,
	token *tokenrepo.TokenRepository,
	loginCodeManager *coderepo.CodeManager,
	oauthApiMap map[authaccount.Provider]oauthapi.OAuthAPI,
) *LoginUsecase {
	return &LoginUsecase{
		authAccount:      authAccount,
		token:            token,
		loginCodeManager: loginCodeManager,
		oauthApiMap:      oauthApiMap,
	}
}
