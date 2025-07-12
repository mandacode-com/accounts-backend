package userevent

import (
	"context"

	"github.com/google/uuid"
	dbrepo "mandacode.com/accounts/auth/internal/repository/database"
)

type UserEventUsecase struct {
	authAccountRepo *dbrepo.AuthAccountRepository
}

func (u *UserEventUsecase) HandleUserDeleted(ctx context.Context, userID uuid.UUID) error {
	if err := u.authAccountRepo.DeleteAuthAccountByUserID(ctx, userID); err != nil {
		return err
	}
	return nil
}

func NewUserEventUsecase(authAccountRepo *dbrepo.AuthAccountRepository) *UserEventUsecase {
	return &UserEventUsecase{
		authAccountRepo: authAccountRepo,
	}
}
