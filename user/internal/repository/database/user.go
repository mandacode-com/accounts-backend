package dbrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/user/ent"
	"mandacode.com/accounts/user/ent/user"
	usermodels "mandacode.com/accounts/user/internal/models/user"
)

type UserRepository struct {
	client *ent.Client
}

// NewUserRepository creates a new UserRepository with the provided database connection string.
func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

// GetUserByID retrieves a user by their ID.
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
	user, err := r.client.User.Query().
		Where(user.IDEQ(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("User not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to get User by ID", errcode.ErrInternalFailure)
	}
	return usermodels.NewSecureUser(user), nil
}

// CreateUser creates a new user with the provided details.
func (r *UserRepository) CreateUser(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
	create := r.client.User.Create().
		SetID(id)

	user, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, errors.New("User already exists", "Conflict", errcode.ErrConflict)
		}
		return nil, errors.New(err.Error(), "Failed to create User", errcode.ErrInternalFailure)
	}
	return usermodels.NewSecureUser(user), nil
}

// UpdateIsActive updates the active status of a user.
func (r *UserRepository) UpdateIsActive(ctx context.Context, id uuid.UUID, isActive bool) (*usermodels.SecureUser, error) {
	user, err := r.client.User.UpdateOneID(id).
		SetIsActive(isActive).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("User not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to update User isActive status", errcode.ErrInternalFailure)
	}
	return usermodels.NewSecureUser(user), nil
}

// DeleteUser deletes a user by their ID.
func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	del := r.client.User.DeleteOneID(id)

	err := del.Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("User not found", "NotFound", errcode.ErrNotFound)
		}
		return errors.New(err.Error(), "Failed to delete User", errcode.ErrInternalFailure)
	}
	return nil
}

// ArchiveUser archives a user by their ID.
func (r *UserRepository) ArchiveUser(ctx context.Context, id uuid.UUID, duration time.Duration, syncCode string) (*usermodels.SecureUser, error) {
	// In this example, archiving is simply setting the isActive field to false.
	user, err := r.client.User.UpdateOneID(id).
		SetIsArchived(true).
		SetArchivedAt(time.Now()).
		SetDeleteAfter(time.Now().Add(duration)).
		SetSyncCode(syncCode).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("User not found", "NotFound", errcode.ErrNotFound)

		}
		return nil, errors.New("Failed to archive User", "InternalError", errcode.ErrInternalFailure)
	}
	return usermodels.NewSecureUser(user), nil
}

// RestoreUser restores a user by their ID.
func (r *UserRepository) RestoreUser(ctx context.Context, id uuid.UUID, syncCode string) (*usermodels.SecureUser, error) {
	user, err := r.client.User.UpdateOneID(id).
		SetIsArchived(false).
		SetNillableArchivedAt(nil).
		SetNillableDeleteAfter(nil).
		SetSyncCode(syncCode).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("User not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to restore User", errcode.ErrInternalFailure)
	}
	return usermodels.NewSecureUser(user), nil
}

// BlockUser blocks a user by their ID.
func (r *UserRepository) BlockUser(ctx context.Context, id uuid.UUID, isBlocked bool, syncCode string) (*usermodels.SecureUser, error) {
	user, err := r.client.User.UpdateOneID(id).
		SetIsBlocked(isBlocked).
		SetSyncCode(syncCode).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("User not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to block User", errcode.ErrInternalFailure)
	}
	return usermodels.NewSecureUser(user), nil
}
