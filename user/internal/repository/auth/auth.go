package authrepo

import (
	"context"

	authv1 "github.com/mandacode-com/accounts-proto/go/auth/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	authrepodto "mandacode.com/accounts/user/internal/repository/auth/dto"
)

type AuthRepository struct {
	localUserClient authv1.LocalUserServiceClient
	oauthUserClient authv1.OAuthUserServiceClient
}

// NewAuthRepository creates a new instance of AuthRepository with the provided clients.
func NewAuthRepository(localUserClient authv1.LocalUserServiceClient, oauthUserClient authv1.OAuthUserServiceClient) *AuthRepository {
	return &AuthRepository{
		localUserClient: localUserClient,
		oauthUserClient: oauthUserClient,
	}
}

func (a *AuthRepository) CreateLocalUser(ctx context.Context, req *authrepodto.CreateLocalUserRequest) (*authrepodto.CreateLocalUserResponse, error) {
	protoRes, err := a.localUserClient.CreateLocalUser(ctx, req.ToProto())
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to create local user", errcode.ErrInternalFailure)
	}
	if err := protoRes.ValidateAll(); err != nil {
		return nil, errors.Upgrade(err, "Invalid response from local user service", errcode.ErrInternalFailure)
	}

	res, err := authrepodto.NewCreateLocalUserResponse(protoRes)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse response from local user service", errcode.ErrInternalFailure)
	}
	return res, nil
}

func (a *AuthRepository) CreateOAuthUser(ctx context.Context, req *authrepodto.CreateOAuthUserRequest) (*authrepodto.CreateOAuthUserResponse, error) {
	protoRes, err := a.oauthUserClient.CreateOAuthUser(ctx, req.ToProto())
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to create OAuth user", errcode.ErrInternalFailure)
	}
	if err := protoRes.ValidateAll(); err != nil {
		return nil, errors.Upgrade(err, "Invalid response from OAuth user service", errcode.ErrInternalFailure)
	}

	res, err := authrepodto.NewCreateOAuthUserResponse(protoRes)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse response from OAuth user service", errcode.ErrInternalFailure)
	}
	return res, nil
}

func (a *AuthRepository) DeleteLocalUser(ctx context.Context, req *authrepodto.DeleteLocalUserRequest) (*authrepodto.DeleteLocalUserResponse, error) {
	protoRes, err := a.localUserClient.DeleteLocalUser(ctx, req.ToProto())
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to delete local user", errcode.ErrInternalFailure)
	}
	if err := protoRes.ValidateAll(); err != nil {
		return nil, errors.Upgrade(err, "Invalid response from local user service", errcode.ErrInternalFailure)
	}

	res, err := authrepodto.NewDeleteLocalUserResponse(protoRes)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse response from local user service", errcode.ErrInternalFailure)
	}
	return res, nil
}

func (a *AuthRepository) DeleteOAuthUser(ctx context.Context, req *authrepodto.DeleteOAuthUserRequest) (*authrepodto.DeleteOAuthUserResponse, error) {
	protoRes, err := a.oauthUserClient.DeleteOAuthUser(ctx, req.ToProto())
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to delete OAuth user", errcode.ErrInternalFailure)
	}
	if err := protoRes.ValidateAll(); err != nil {
		return nil, errors.Upgrade(err, "Invalid response from OAuth user service", errcode.ErrInternalFailure)
	}

	res, err := authrepodto.NewDeleteOAuthUserResponse(protoRes)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse response from OAuth user service", errcode.ErrInternalFailure)
	}
	return res, nil
}

func (a *AuthRepository) UpdateLocalUserEmail(ctx context.Context, req *authrepodto.UpdateLocalUserEmailRequest) (*authrepodto.UpdateLocalUserEmailResponse, error) {
	protoRes, err := a.localUserClient.UpdateLocalUserEmail(ctx, req.ToProto())
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to update local user email", errcode.ErrInternalFailure)
	}
	if err := protoRes.ValidateAll(); err != nil {
		return nil, errors.Upgrade(err, "Invalid response from local user service", errcode.ErrInternalFailure)
	}

	res, err := authrepodto.NewUpdateLocalUserEmailResponse(protoRes)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse response from local user service", errcode.ErrInternalFailure)
	}
	return res, nil
}

func (a *AuthRepository) UpdateLocalUserEmailVerification(ctx context.Context, req *authrepodto.UpdateEmailVerificationRequest) (*authrepodto.UpdateEmailVerificationResponse, error) {
	protoRes, err := a.localUserClient.UpdateEmailVerification(ctx, req.ToProto())
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to update email verification", errcode.ErrInternalFailure)
	}
	if err := protoRes.ValidateAll(); err != nil {
		return nil, errors.Upgrade(err, "Invalid response from local user service", errcode.ErrInternalFailure)
	}

	res, err := authrepodto.NewUpdateEmailVerificationResponse(protoRes)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse response from local user service", errcode.ErrInternalFailure)
	}
	return res, nil
}

func (a *AuthRepository) SyncOAuthUser(ctx context.Context, req *authrepodto.SyncOAuthUserRequest) (*authrepodto.SyncOAuthUserResponse, error) {
	protoRes, err := a.oauthUserClient.SyncOAuthUser(ctx, req.ToProto())
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to sync OAuth user", errcode.ErrInternalFailure)
	}
	if err := protoRes.ValidateAll(); err != nil {
		return nil, errors.Upgrade(err, "Invalid response from OAuth user service", errcode.ErrInternalFailure)
	}

	res, err := authrepodto.NewSyncOAuthUserResponse(protoRes)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse response from OAuth user service", errcode.ErrInternalFailure)
	}
	return res, nil
}
