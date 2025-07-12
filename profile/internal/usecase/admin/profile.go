package admin

import (
	"context"

	"github.com/google/uuid"
	dbmodels "mandacode.com/accounts/profile/internal/models/database"
	dbrepo "mandacode.com/accounts/profile/internal/repository/database"
	"mandacode.com/accounts/profile/internal/usecase/dto"
)

type ProfileUsecase struct {
	repo *dbrepo.ProfileRepository
}

func NewProfileUsecase(repo *dbrepo.ProfileRepository) *ProfileUsecase {
	return &ProfileUsecase{
		repo: repo,
	}
}

func (u *ProfileUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*dbmodels.SecureProfile, error) {
	prof, err := u.repo.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dbmodels.NewSecureProfile(prof), nil
}

func (u *ProfileUsecase) UpdateProfile(ctx context.Context, data *dto.UpdateProfileData) (*dbmodels.SecureProfile, error) {
	prof, err := u.repo.UpdateProfile(ctx, data.ToRepoModel())
	if err != nil {
		return nil, err
	}

	return dbmodels.NewSecureProfile(prof), nil
}
