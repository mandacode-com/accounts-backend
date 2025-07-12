package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	usermodels "mandacode.com/accounts/user/internal/models/user"
	dbrepo "mandacode.com/accounts/user/internal/repository/database"
	usereventrepo "mandacode.com/accounts/user/internal/repository/userevent"
	"mandacode.com/accounts/user/internal/util"
)

type UserUsecase struct {
	userRepo      *dbrepo.UserRepository
	eventEmitter  *usereventrepo.UserEventEmitter
	deleteDelay   time.Duration
	codeGenerator *util.RandomStringGenerator
}

// NewUserUsecase creates a new UserUsecase with the provided repositories.
func NewUserUsecase(userRepo *dbrepo.UserRepository, eventEmitter *usereventrepo.UserEventEmitter) *UserUsecase {
	return &UserUsecase{
		userRepo:     userRepo,
		eventEmitter: eventEmitter,
	}
}

// GetUserByID retrieves a user by their ID.
func (u *UserUsecase) GetUserByID(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user with the provided details.
func (u *UserUsecase) CreateUser(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
	user, err := u.userRepo.CreateUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateIsActive updates the active status of a user.
func (u *UserUsecase) UpdateIsActive(ctx context.Context, id uuid.UUID, isActive bool) (*usermodels.SecureUser, error) {
	user, err := u.userRepo.UpdateIsActive(ctx, id, isActive)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by their ID.
func (u *UserUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := u.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	// Emit a user deletion event
	if err := u.eventEmitter.EmitUserDeletedEvent(ctx, id); err != nil {
		return err
	}

	return nil
}

// ArchiveUser archives a user by their ID.
func (u *UserUsecase) ArchiveUser(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
	syncCode, err := u.codeGenerator.Generate()
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.ArchiveUser(ctx, id, u.deleteDelay, syncCode)
	if err != nil {
		return nil, err
	}

	// Emit a user archived event
	if err := u.eventEmitter.EmitUserArchivedEvent(ctx, id, syncCode); err != nil {
		return nil, err
	}

	return user, nil
}

// RestoreUser restores a user by their ID.
func (u *UserUsecase) RestoreUser(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
	syncCode, err := u.codeGenerator.Generate()
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.RestoreUser(ctx, id, syncCode)
	if err != nil {
		return nil, err
	}

	// Emit a user restored event
	if err := u.eventEmitter.EmitUserRestoredEvent(ctx, id, syncCode); err != nil {
		return nil, err
	}

	return user, nil
}

// BlockUser blocks a user by their ID.
func (u *UserUsecase) BlockUser(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
	syncCode, err := u.codeGenerator.Generate()
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.BlockUser(ctx, id, true, syncCode)
	if err != nil {
		return nil, err
	}
	// Emit a user blocked event
	if err := u.eventEmitter.EmitUserBlockedEvent(ctx, id, syncCode); err != nil {
		return nil, err
	}
	return user, nil
}

// UnblockUser unblocks a user by their ID.
func (u *UserUsecase) UnblockUser(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
	syncCode, err := u.codeGenerator.Generate()
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.BlockUser(ctx, id, false, syncCode)
	if err != nil {
		return nil, err
	}

	// Emit a user unblocked event
	if err := u.eventEmitter.EmitUserUnblockedEvent(ctx, id, syncCode); err != nil {
		return nil, err
	}

	return user, nil
}
