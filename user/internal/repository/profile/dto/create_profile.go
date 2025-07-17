package profilerepodto

import (
	"time"

	"github.com/google/uuid"
	profilev1 "github.com/mandacode-com/accounts-proto/go/profile/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
)

type CreateProfileUserRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	SyncCode string    `json:"sync_code"`
}

func (r *CreateProfileUserRequest) ToProto() *profilev1.InitUserRequest {
	if r == nil {
		return nil
	}
	return &profilev1.InitUserRequest{
		UserId:   r.UserID.String(),
		Email:    r.Email,
		SyncCode: r.SyncCode,
	}
}

type CreateProfileUserResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt string    `json:"created_at"`
}

func NewCreateProfileUserResponse(data *profilev1.InitUserResponse) (*CreateProfileUserResponse, error) {
	if data == nil {
		return nil, errors.New("data is nil", "Invalid response data", errcode.ErrInternalFailure)
	}
	userUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse user ID", errcode.ErrInvalidFormat)
	}
	return &CreateProfileUserResponse{
		UserID:    userUID,
		CreatedAt: data.InitializedAt.AsTime().Format(time.RFC3339),
	}, nil
}
