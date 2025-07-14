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
	CreateLocalAuthUser(ctx context.Context, userID uuid.UUID, email string, password string) (uuid.UUID, error)
	DeleteAuthUser(ctx context.Context, userID uuid.UUID) error
	UpdateAuthUserEmail(ctx context.Context, userID uuid.UUID, newEmail string) error
}

type localUserUsecase struct {
	authAccountRepo *dbrepo.AuthAccountRepository
}

// CreateLocalAuthUser implements IAuthUserUsecase.
func (a *localUserUsecase) CreateLocalAuthUser(ctx context.Context, userID uuid.UUID, email string, password string) (uuid.UUID, error) {
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
		return uuid.Nil, errors.Upgrade(err, "Failed to create local auth account", errcode.ErrInternalFailure)
	}
	return account.UserID, nil
}

// DeleteAuthUser implements IAuthUserUsecase.
func (a *localUserUsecase) DeleteAuthUser(ctx context.Context, userID uuid.UUID) error {
	if err := a.authAccountRepo.DeleteAuthAccountByUserID(ctx, userID); err != nil {
		return err
	}
	return nil
}

// UpdateLocalAuthUserEmail implements IAuthUserUsecase.
func (a *localUserUsecase) UpdateAuthUserEmail(ctx context.Context, userID uuid.UUID, newEmail string) error {
	account, err := a.authAccountRepo.GetLocalAuthAccountByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if _, err = a.authAccountRepo.UpdateEmailByID(ctx, account.ID, newEmail); err != nil {
		return err
	}

	return nil
}

func NewLocalUserUsecase(authAccountRepo *dbrepo.AuthAccountRepository) LocalUserUsecase {
	return &localUserUsecase{
		authAccountRepo: authAccountRepo,
	}
}
