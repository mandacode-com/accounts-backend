package login

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/auth/ent/authaccount"

	"mandacode.com/accounts/auth/internal/infra/oauthapi"
	signupinfra "mandacode.com/accounts/auth/internal/infra/signup"
	coderepo "mandacode.com/accounts/auth/internal/repository/code"
	dbrepo "mandacode.com/accounts/auth/internal/repository/database"
	tokenrepo "mandacode.com/accounts/auth/internal/repository/token"
	logindto "mandacode.com/accounts/auth/internal/usecase/login/dto"
)

type OAuthLoginUsecase struct {
	authAccount      *dbrepo.AuthAccountRepository
	token            *tokenrepo.TokenRepository
	loginCodeManager *coderepo.CodeManager
	signupApi        *signupinfra.SignupAPI
	oauthApiMap      map[authaccount.Provider]oauthapi.OAuthAPI
}

// getAccessToken retrieves the access token from the OAuth API.
func (l *OAuthLoginUsecase) getAccessToken(ctx context.Context, provider authaccount.Provider, code string) (string, error) {
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

// getOrCreateVerifiedUser retrieves or creates a verified user based on the OAuth input.
func (l *OAuthLoginUsecase) getOrCreateVerifiedUser(ctx context.Context, input logindto.OAuthLoginInput) (uuid.UUID, error) {
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
		if errors.Is(err, errcode.ErrNotFound) { // If the OAuth account does not exist, create a new one
			signupResponse, err := l.signupApi.OAuthSignup(input.Provider, oauthAccessToken) // Request signup API to create a new user
			if err != nil {
				return uuid.Nil, err
			}
			userUID, err := uuid.Parse(signupResponse.UserID)
			if err != nil {
				return uuid.Nil, errors.Upgrade(err, "Failed to parse user ID from signup response", errcode.ErrInternalFailure)
			}
			userID = userUID
			verified = signupResponse.IsVerified
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

// GetLoginURL implements oauthdomain.OAuthLoginUsecase.
func (l *OAuthLoginUsecase) GetLoginURL(ctx context.Context, provider string) (loginURL string, err error) {
	api, ok := l.oauthApiMap[authaccount.Provider(provider)]
	if !ok {
		return "", errors.New("unsupported provider: "+provider, "Unsupported Provider", errcode.ErrInvalidInput)
	}
	loginURL = api.GetLoginURL()
	return loginURL, nil
}

// IssueLoginCode implements oauthdomain.OAuthLoginUsecase.
func (l *OAuthLoginUsecase) IssueLoginCode(ctx context.Context, input logindto.OAuthLoginInput) (code string, userID uuid.UUID, err error) {
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

// Login implements oauthdomain.OAuthLoginUsecase.
func (l *OAuthLoginUsecase) Login(ctx context.Context, input logindto.OAuthLoginInput) (accessToken string, refreshToken string, err error) {
	// Get or create verified user
	userID, err := l.getOrCreateVerifiedUser(ctx, input)
	if err != nil {
		return "", "", errors.Upgrade(err, "Failed to get or create verified user", errcode.ErrUnauthorized)
	}

	// Generate access and refresh tokens
	return l.issueToken(ctx, userID)
}

// VerifyLoginCode implements oauthdomain.OAuthLoginUsecase.
func (l *OAuthLoginUsecase) VerifyLoginCode(ctx context.Context, userID uuid.UUID, code string) (accessToken string, refreshToken string, err error) {
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
func (l *OAuthLoginUsecase) issueToken(ctx context.Context, userID uuid.UUID) (accessToken string, refreshToken string, err error) {
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

// NewOAuthLoginUsecase creates a new instance of LoginUsecase.
func NewOAuthLoginUsecase(
	authAccount *dbrepo.AuthAccountRepository,
	token *tokenrepo.TokenRepository,
	loginCodeManager *coderepo.CodeManager,
	signupApi *signupinfra.SignupAPI,
	oauthApiMap map[authaccount.Provider]oauthapi.OAuthAPI,
) *OAuthLoginUsecase {
	return &OAuthLoginUsecase{
		authAccount:      authAccount,
		token:            token,
		loginCodeManager: loginCodeManager,
		signupApi:        signupApi,
		oauthApiMap:      oauthApiMap,
	}
}
