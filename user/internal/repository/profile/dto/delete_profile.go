package profilerepodto

import (
	"github.com/google/uuid"
	profilev1 "github.com/mandacode-com/accounts-proto/go/profile/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"time"
)

type DeleteProfileUserRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

func (r *DeleteProfileUserRequest) ToProto() *profilev1.DeleteUserRequest {
	if r == nil {
		return nil
	}
	return &profilev1.DeleteUserRequest{
		UserId: r.UserID.String(),
	}
}

type DeleteProfileUserResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	DeletedAt string    `json:"deleted_at"`
}

func NewDeleteProfileUserResponse(data *profilev1.DeleteUserResponse) (*DeleteProfileUserResponse, error) {
	if data == nil {
		return nil, errors.New("data is nil", "Invalid response data", errcode.ErrInternalFailure)
	}
	userUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse user ID", errcode.ErrInvalidFormat)
	}
	return &DeleteProfileUserResponse{
		UserID:    userUID,
		DeletedAt: data.DeletedAt.AsTime().Format(time.RFC3339),
	}, nil
}
