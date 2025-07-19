package authrepodto

import (
	"time"

	"github.com/google/uuid"
	authv1 "github.com/mandacode-com/accounts-proto/go/auth/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/user/internal/models/provider"
)

type CreateLocalUserRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func (r *CreateLocalUserRequest) ToProto() *authv1.CreateLocalUserRequest {
	if r == nil {
		return nil
	}
	return &authv1.CreateLocalUserRequest{
		UserId:   r.UserID.String(),
		Email:    r.Email,
		Password: r.Password,
	}
}

type CreateLocalUserResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewCreateLocalUserResponse(data *authv1.CreateLocalUserResponse) (*CreateLocalUserResponse, error) {
	if data == nil {
		return nil, errors.New("data is nil", "Invalid response data", errcode.ErrInternalFailure)
	}
	userUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse user ID", errcode.ErrInvalidFormat)
	}
	return &CreateLocalUserResponse{
		UserID:    userUID,
		CreatedAt: data.CreatedAt.AsTime(),
	}, nil
}

type CreateOAuthUserRequest struct {
	UserID      uuid.UUID             `json:"user_id"`
	Provider    provider.ProviderType `json:"provider"`
	AccessToken *string               `json:"access_token"`
	Code        *string               `json:"code"`
}

func (r *CreateOAuthUserRequest) ToProto() *authv1.CreateOAuthUserRequest {
	if r == nil {
		return nil
	}
	return &authv1.CreateOAuthUserRequest{
		UserId:      r.UserID.String(),
		Provider:    r.Provider.ToProto(),
		AccessToken: r.AccessToken,
		Code:        r.Code,
	}
}

type CreateOAuthUserResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	Provider  string    `json:"provider"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewCreateOAuthUserResponse(data *authv1.CreateOAuthUserResponse) (*CreateOAuthUserResponse, error) {
	if data == nil {
		return nil, errors.New("data is nil", "Invalid response data", errcode.ErrInternalFailure)
	}
	userUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse user ID", errcode.ErrInvalidFormat)
	}
	return &CreateOAuthUserResponse{
		UserID:    userUID,
		Provider:  data.Provider.String(),
		Email:     data.Email,
		CreatedAt: data.CreatedAt.AsTime(),
	}, nil
}
