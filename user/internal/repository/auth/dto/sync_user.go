package authrepodto

import (
	"time"

	"github.com/google/uuid"
	authv1 "github.com/mandacode-com/accounts-proto/go/auth/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/user/internal/models/provider"
)

type SyncOAuthUserRequest struct {
	UserID      uuid.UUID             `json:"user_id"`
	Provider    provider.ProviderType `json:"provider"`
	AccessToken *string               `json:"access_token"`
	Code        *string               `json:"code"`
}

func (r *SyncOAuthUserRequest) ToProto() *authv1.SyncOAuthUserRequest {
	if r == nil {
		return nil
	}
	return &authv1.SyncOAuthUserRequest{
		UserId:      r.UserID.String(),
		Provider:    r.Provider.ToProto(),
		AccessToken: r.AccessToken,
		Code:        r.Code,
	}
}

type SyncOAuthUserResponse struct {
	UserID   uuid.UUID             `json:"user_id"`
	Provider provider.ProviderType `json:"provider"`
	SyncedAt time.Time             `json:"synced_at"`
}

func NewSyncOAuthUserResponse(data *authv1.SyncOAuthUserResponse) (*SyncOAuthUserResponse, error) {
	if data == nil {
		return nil, errors.New("data is nil", "Invalid response data", errcode.ErrInternalFailure)
	}
	userUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse user ID", errcode.ErrInvalidFormat)
	}
	localProvider := provider.ToLocalProvider(data.Provider)
	return &SyncOAuthUserResponse{
		UserID:   userUID,
		Provider: localProvider,
		SyncedAt: data.SyncedAt.AsTime(),
	}, nil
}
