package profilerepodto

import (
	"time"

	"github.com/google/uuid"
	profilev1 "github.com/mandacode-com/accounts-proto/go/profile/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
)

type UpdateEmailRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	SyncCode string    `json:"sync_code"`
}

func (r *UpdateEmailRequest) ToProto() *profilev1.UpdateEmailRequest {
	if r == nil {
		return nil
	}
	return &profilev1.UpdateEmailRequest{
		UserId:   r.UserID.String(),
		NewEmail: r.Email,
	}
}

type UpdateEmailResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	UpdatedAt string    `json:"updated_at"`
}

func NewUpdateEmailResponse(data *profilev1.UpdateEmailResponse) (*UpdateEmailResponse, error) {
	if data == nil {
		return nil, errors.New("data is nil", "Invalid response data", errcode.ErrInternalFailure)
	}
	userUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse user ID", errcode.ErrInvalidFormat)
	}
	return &UpdateEmailResponse{
		UserID:    userUID,
		Email:     data.UpdatedEmail,
		UpdatedAt: data.UpdatedAt.AsTime().Format(time.RFC3339),
	}, nil
}
