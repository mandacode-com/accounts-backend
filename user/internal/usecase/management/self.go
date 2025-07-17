package manage

import (
	"context"
	"time"

	"github.com/google/uuid"
	usermodels "mandacode.com/accounts/user/internal/models/user"
	dbrepo "mandacode.com/accounts/user/internal/repository/database"
	usereventrepo "mandacode.com/accounts/user/internal/repository/userevent"
	"mandacode.com/accounts/user/internal/util"
)

type SelfManageUsecase struct {
	userRepo      *dbrepo.UserRepository
	eventEmitter  *usereventrepo.UserEventEmitter
	deleteDelay   time.Duration
	codeGenerator *util.RandomStringGenerator
}

// NewManageUsecase creates a new ManageUsecase with the provided repositories.
func NewSelfManageUsecase(userRepo *dbrepo.UserRepository, eventEmitter *usereventrepo.UserEventEmitter) *SelfManageUsecase {
	return &SelfManageUsecase{
		userRepo:      userRepo,
		eventEmitter:  eventEmitter,
		deleteDelay:   24 * time.Hour,                    // Default delete delay of 24 hours
		codeGenerator: util.NewRandomStringGenerator(32), // Default code length of 32 characters
	}
}

// GetUserByID retrieves a user by their ID.
func (m *SelfManageUsecase) GetUserByID(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
	user, err := m.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ArchiveUser archives a user by their ID.
func (m *SelfManageUsecase) ArchiveUser(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
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
func (m *SelfManageUsecase) RestoreUser(ctx context.Context, id uuid.UUID) (*usermodels.SecureUser, error) {
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
