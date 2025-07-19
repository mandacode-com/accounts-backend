package manage

import (
	"context"
	"time"

	"github.com/google/uuid"
	usermodels "mandacode.com/accounts/user/internal/models/user"
	dbrepo "mandacode.com/accounts/user/internal/repository/database"
	usereventrepo "mandacode.com/accounts/user/internal/repository/userevent"
)

type AdminManageUsecase struct {
	userRepo     *dbrepo.UserRepository
	eventEmitter *usereventrepo.UserEventEmitter
	deleteDelay  time.Duration
}

// NewAdminManageUsecase creates a new ManageUsecase with the provided repositories.
func NewAdminManageUsecase(userRepo *dbrepo.UserRepository, eventEmitter *usereventrepo.UserEventEmitter) *AdminManageUsecase {
	return &AdminManageUsecase{
		userRepo:     userRepo,
		eventEmitter: eventEmitter,
		deleteDelay:  24 * time.Hour, // Default delete delay of 24 hours
	}
}

// GetUserByID retrieves a user by their ID.
func (m *AdminManageUsecase) GetUserByID(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
	user, err := m.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ArchiveUser archives a user by their ID.
func (m *AdminManageUsecase) ArchiveUser(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
	user, err := m.userRepo.ArchiveUser(ctx, id, m.deleteDelay)
	if err != nil {
		return nil, err
	}

	// Emit an event for archiving the user
	if err := m.eventEmitter.EmitUserArchivedEvent(ctx, user.ID, user.SyncCode); err != nil {
		return nil, err
	}

	return user, nil
}

// RestoreUser restores an archived user by their ID.
func (m *AdminManageUsecase) RestoreUser(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
	user, err := m.userRepo.RestoreUser(ctx, id)
	if err != nil {
		return nil, err
	}

	// Emit an event for restoring the user
	if err := m.eventEmitter.EmitUserRestoredEvent(ctx, user.ID, user.SyncCode); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by their ID.
func (m *AdminManageUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := m.userRepo.DeleteUser(ctx, id); err != nil {
		return err
	}

	// Emit a user deletion event
	if err := m.eventEmitter.EmitUserDeletedEvent(ctx, id); err != nil {
		return err
	}

	return nil
}

// BlockUser blocks a user by their ID.
func (m *AdminManageUsecase) BlockUser(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
	user, err := m.userRepo.BlockUser(ctx, id, true)
	if err != nil {
		return nil, err
	}

	// Emit a user blocked event
	if err := m.eventEmitter.EmitUserBlockedEvent(ctx, user.ID, user.SyncCode); err != nil {
		return nil, err
	}

	return user, nil
}

// UnblockUser unblocks a user by their ID.
func (m *AdminManageUsecase) UnblockUser(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
	user, err := m.userRepo.BlockUser(ctx, id, false)
	if err != nil {
		return nil, err
	}

	// Emit a user unblocked event
	if err := m.eventEmitter.EmitUserUnblockedEvent(ctx, user.ID, user.SyncCode); err != nil {
		return nil, err
	}

	return user, nil
}
