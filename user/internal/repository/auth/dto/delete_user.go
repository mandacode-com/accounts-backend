package authrepodto

import (
	"time"

	"github.com/google/uuid"
	authv1 "github.com/mandacode-com/accounts-proto/go/auth/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/user/internal/models/provider"
)

type DeleteLocalUserRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

func (r *DeleteLocalUserRequest) ToProto() *authv1.DeleteLocalUserRequest {
	if r == nil {
		return nil
	}
	return &authv1.DeleteLocalUserRequest{
		UserId: r.UserID.String(),
	}
}

type DeleteLocalUserResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	DeletedAt time.Time `json:"deleted_at"`
}

func NewDeleteLocalUserResponse(data *authv1.DeleteLocalUserResponse) (*DeleteLocalUserResponse, error) {
	if data == nil {
		return nil, errors.New("data is nil", "Invalid response data", errcode.ErrInternalFailure)
	}
	userUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse user ID", errcode.ErrInvalidFormat)
	}
	return &DeleteLocalUserResponse{
		UserID:    userUID,
		DeletedAt: data.DeletedAt.AsTime(),
	}, nil
}

type DeleteOAuthUserRequest struct {
	UserID   uuid.UUID             `json:"user_id"`
	Provider provider.ProviderType `json:"provider"`
}

func (r *DeleteOAuthUserRequest) ToProto() *authv1.DeleteOAuthUserRequest {
	if r == nil {
		return nil
	}
	return &authv1.DeleteOAuthUserRequest{
		UserId:   r.UserID.String(),
		Provider: r.Provider.ToProto(),
	}
}

type DeleteOAuthUserResponse struct {
	UserID    uuid.UUID             `json:"user_id"`
	Provider  provider.ProviderType `json:"provider"`
	DeletedAt time.Time             `json:"deleted_at"`
}

func NewDeleteOAuthUserResponse(data *authv1.DeleteOAuthUserResponse) (*DeleteOAuthUserResponse, error) {
	if data == nil {
		return nil, errors.New("data is nil", "Invalid response data", errcode.ErrInternalFailure)
	}
	userUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse user ID", errcode.ErrInvalidFormat)
	}
	localProvider := provider.ToLocalProvider(data.Provider)
	return &DeleteOAuthUserResponse{
		UserID:    userUID,
		Provider:  localProvider,
		DeletedAt: data.DeletedAt.AsTime(),
	}, nil
}
