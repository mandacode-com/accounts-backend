package authrepodto

import (
	"time"

	"github.com/google/uuid"
	authv1 "github.com/mandacode-com/accounts-proto/go/auth/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
)

type UpdateLocalUserEmailRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	NewEmail string    `json:"new_email"`
	SyncCode string    `json:"sync_code"`
}

func (r *UpdateLocalUserEmailRequest) ToProto() *authv1.UpdateLocalUserEmailRequest {
	if r == nil {
		return nil
	}
	return &authv1.UpdateLocalUserEmailRequest{
		UserId:   r.UserID.String(),
		NewEmail: r.NewEmail,
	}
}

type UpdateLocalUserEmailResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	UpdatedEmail string    `json:"updated_email"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewUpdateLocalUserEmailResponse(data *authv1.UpdateLocalUserEmailResponse) (*UpdateLocalUserEmailResponse, error) {
	if data == nil {
		return nil, errors.New("data is nil", "Invalid response data", errcode.ErrInternalFailure)
	}
	userUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse user ID", errcode.ErrInvalidFormat)
	}
	return &UpdateLocalUserEmailResponse{
		UserID:       userUID,
		UpdatedEmail: data.UpdatedEmail,
		UpdatedAt:    data.UpdatedAt.AsTime(),
	}, nil
}

type UpdateEmailVerificationRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	Verified bool      `json:"verified"`
}

func (r *UpdateEmailVerificationRequest) ToProto() *authv1.UpdateEmailVerificationRequest {
	if r == nil {
		return nil
	}
	return &authv1.UpdateEmailVerificationRequest{
		UserId:   r.UserID.String(),
		Verified: r.Verified,
	}
}

type UpdateEmailVerificationResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	Verified  bool      `json:"verified"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUpdateEmailVerificationResponse(data *authv1.UpdateEmailVerificationResponse) (*UpdateEmailVerificationResponse, error) {
	if data == nil {
		return nil, errors.New("data is nil", "Invalid response data", errcode.ErrInternalFailure)
	}
	userUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse user ID", errcode.ErrInvalidFormat)
	}
	return &UpdateEmailVerificationResponse{
		UserID:    userUID,
		Verified:  data.Verified,
		UpdatedAt: data.UpdatedAt.AsTime(),
	}, nil
}
