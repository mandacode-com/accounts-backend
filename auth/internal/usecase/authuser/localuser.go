package authuser

import (
	"context"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	dbmodels "mandacode.com/accounts/auth/internal/models/database"
	dbrepo "mandacode.com/accounts/auth/internal/repository/database"
)

type LocalUserUsecase interface {
	CreateLocalAuthUser(ctx context.Context, userID uuid.UUID, email string, password string) (*dbmodels.SecureLocalAuthAccount, error)
	DeleteAuthUser(ctx context.Context, userID uuid.UUID) error
	UpdateAuthUserEmail(ctx context.Context, userID uuid.UUID, newEmail string) (*dbmodels.SecureAuthAccount, error)
	UpdateLocalEmailVerificationStatus(ctx context.Context, userID uuid.UUID, isVerified bool) error
}

type localUserUsecase struct {
	authAccountRepo *dbrepo.AuthAccountRepository
}

// CreateLocalAuthUser implements IAuthUserUsecase.
func (a *localUserUsecase) CreateLocalAuthUser(ctx context.Context, userID uuid.UUID, email string, password string) (*dbmodels.SecureLocalAuthAccount, error) {
	account, err := a.authAccountRepo.CreateLocalAuthAccount(
		ctx,
		&dbmodels.CreateLocalAuthAccountInput{
			UserID:     userID,
			Email:      email,
			Password:   password,
			IsVerified: false, // Default to false, verification can be handled later
		},
	)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to create local auth account", errcode.ErrInternalFailure)
	}
	return account, nil
}

// DeleteAuthUser implements IAuthUserUsecase.
func (a *localUserUsecase) DeleteAuthUser(ctx context.Context, userID uuid.UUID) error {
	if err := a.authAccountRepo.DeleteAuthAccountByUserID(ctx, userID); err != nil {
		return err
	}
	return nil
}

// UpdateLocalAuthUserEmail implements IAuthUserUsecase.
func (a *localUserUsecase) UpdateAuthUserEmail(ctx context.Context, userID uuid.UUID, newEmail string) (*dbmodels.SecureAuthAccount, error) {
	account, err := a.authAccountRepo.GetLocalAuthAccountByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to get local auth account by user ID", errcode.ErrInternalFailure)
	}
	updatedAccount, err := a.authAccountRepo.UpdateEmailByID(ctx, account.ID, newEmail)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to update local auth account email", errcode.ErrInternalFailure)
	}

	return updatedAccount, nil
}

// UpdateEmailVerificationStatus implements IAuthUserUsecase.
func (a *localUserUsecase) UpdateLocalEmailVerificationStatus(ctx context.Context, userID uuid.UUID, isVerified bool) error {
	if err := a.authAccountRepo.UpdateLocalEmailVerificationStatus(ctx, userID, isVerified); err != nil {
		return errors.Upgrade(err, "Failed to update email verification status", errcode.ErrInternalFailure)
	}

	return nil
}

func NewLocalUserUsecase(authAccountRepo *dbrepo.AuthAccountRepository) LocalUserUsecase {
	return &localUserUsecase{
		authAccountRepo: authAccountRepo,
	}
}
