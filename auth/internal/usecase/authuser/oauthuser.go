package authuser

import (
	"context"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/auth/ent/authaccount"
	"mandacode.com/accounts/auth/internal/infra/oauthapi"
	dbmodels "mandacode.com/accounts/auth/internal/models/database"
	dbrepo "mandacode.com/accounts/auth/internal/repository/database"
)

type OAuthUserUsecase interface {
	CreateOAuthUser(ctx context.Context, userID uuid.UUID, provider authaccount.Provider, accessToken *string, code *string) (*dbmodels.SecureOAuthAuthAccount, error)
	DeleteOAuthUser(ctx context.Context, userID uuid.UUID) error
	SyncOAuthUser(ctx context.Context, userID uuid.UUID, provider authaccount.Provider, accessToken *string, code *string) (*dbmodels.SecureOAuthAuthAccount, error)
}

type oauthUserUsecase struct {
	authAccountRepo *dbrepo.AuthAccountRepository
	oauthApiMap     map[authaccount.Provider]oauthapi.OAuthAPI
}

// CreateOAuthUser implements IAuthUserUsecase.
func (a *oauthUserUsecase) CreateOAuthUser(ctx context.Context, userID uuid.UUID, provider authaccount.Provider, accessToken *string, code *string) (*dbmodels.SecureOAuthAuthAccount, error) {
	api, ok := a.oauthApiMap[provider]
	if !ok {
		return nil, errors.New("unsupported provider: "+string(provider), "UnsupportedProvider", errcode.ErrInvalidInput)
	}

	var oauthAccessToken string
	if accessToken != nil {
		oauthAccessToken = *accessToken
	} else if code != nil {
		var err error
		oauthAccessToken, err = api.GetAccessToken(*code)
		if err != nil {
			return nil, errors.Upgrade(err, "Failed to get access token from OAuth provider", errcode.ErrUnauthorized)
		}
	} else {
		return nil, errors.New("either access token or code must be provided", "InvalidInput", errcode.ErrInvalidInput)
	}

	userInfo, err := api.GetUserInfo(oauthAccessToken)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to get user info from OAuth provider", errcode.ErrUnauthorized)
	}

	account, err := a.authAccountRepo.CreateOAuthAuthAccount(
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

// DeleteOAuthUser implements IAuthUserUsecase.
func (a *oauthUserUsecase) DeleteOAuthUser(ctx context.Context, userID uuid.UUID) error {
	if err := a.authAccountRepo.DeleteAuthAccountByUserID(ctx, userID); err != nil {
		return err
	}
	return nil
}

// SyncOAuthUser implements IAuthUserUsecase.
func (a *oauthUserUsecase) SyncOAuthUser(ctx context.Context, userID uuid.UUID, provider authaccount.Provider, accessToken *string, code *string) (*dbmodels.SecureOAuthAuthAccount, error) {
	api, ok := a.oauthApiMap[provider]
	if !ok {
		return nil, errors.New("unsupported provider: "+string(provider), "UnsupportedProvider", errcode.ErrInvalidInput)
	}

	var oauthAccessToken string
	if accessToken != nil {
		oauthAccessToken = *accessToken
	} else if code != nil {
		var err error
		oauthAccessToken, err = api.GetAccessToken(*code)
		if err != nil {
			return nil, errors.Upgrade(err, "Failed to get access token from OAuth provider", errcode.ErrUnauthorized)
		}
	} else {
		return nil, errors.New("either access token or code must be provided", "InvalidInput", errcode.ErrInvalidInput)
	}

	userInfo, err := api.GetUserInfo(oauthAccessToken)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to get user info from OAuth provider", errcode.ErrUnauthorized)
	}

	account, err := a.authAccountRepo.GetOAuthAccountByProviderAndProviderID(
		ctx,
		provider,
		userInfo.ProviderID,
	)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to sync OAuth account", errcode.ErrInternalFailure)
	}

	_, err = a.authAccountRepo.UpdateEmailByID(
		ctx,
		account.ID,
		userInfo.Email,
	)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to update OAuth account email", errcode.ErrInternalFailure)
	}

	return account, nil
}

func NewOAuthUserUsecase(authAccountRepo *dbrepo.AuthAccountRepository, oauthApis map[authaccount.Provider]oauthapi.OAuthAPI) OAuthUserUsecase {
	return &oauthUserUsecase{
		authAccountRepo: authAccountRepo,
		oauthApiMap:     oauthApis,
	}
}
