package userrepo

import (
	"context"

	"github.com/google/uuid"
	userv1 "github.com/mandacode-com/accounts-proto/go/user/user/v1"
)

type UserServiceRepository struct {
	client userv1.UserServiceClient
}

// InitUser initialize user
func (u *UserServiceRepository) InitUser(ctx context.Context, userID uuid.UUID) (*userv1.InitUserResponse, error) {
	resp, err := u.client.InitUser(ctx, &userv1.InitUserRequest{
		UserId: userID.String(),
	})
	if err != nil {
		return nil, err
	}
	if err := resp.ValidateAll(); err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteUser deletes a user from the system.
func (u *UserServiceRepository) DeleteUser(ctx context.Context, userID uuid.UUID) (*userv1.DeleteUserResponse, error) {
	resp, err := u.client.DeleteUser(ctx, &userv1.DeleteUserRequest{
		UserId: userID.String(),
	})
	if err != nil {
		return nil, err
	}
	if err := resp.ValidateAll(); err != nil {
		return nil, err
	}
	return resp, nil
}

// IsBlocked checks if a user is blocked.
func (u *UserServiceRepository) IsBlocked(ctx context.Context, userID uuid.UUID) (*userv1.IsBlockedResponse, error) {
	resp, err := u.client.IsBlocked(ctx, &userv1.IsBlockedRequest{
		UserId: userID.String(),
	})
	if err != nil {
		return nil, err
	}
	if err := resp.ValidateAll(); err != nil {
		return nil, err
	}
	return resp, nil
}

// IsActive checks if a user is active.
func (u *UserServiceRepository) IsActive(ctx context.Context, userID uuid.UUID) (*userv1.IsActiveResponse, error) {
	resp, err := u.client.IsActive(ctx, &userv1.IsActiveRequest{
		UserId: userID.String(),
	})
	if err != nil {
		return nil, err
	}
	if err := resp.ValidateAll(); err != nil {
		return nil, err
	}
	return resp, nil
}
