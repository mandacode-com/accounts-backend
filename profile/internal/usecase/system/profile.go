package system

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

func (u *ProfileUsecase) CreateProfile(ctx context.Context, data *dto.CreateProfileData) (*dbmodels.SecureProfile, error) {
	prof, err := u.repo.CreateProfile(ctx, data.ToRepoModel())
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

func (u *ProfileUsecase) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	if err := u.repo.DeleteProfile(ctx, userID); err != nil {
		return err
	}
	return nil
}

func (u *ProfileUsecase) ArchiveProfile(ctx context.Context, userID uuid.UUID) error {
	if err := u.repo.ArchiveProfile(ctx, userID); err != nil {
		return err
	}
	return nil
}

func (u *ProfileUsecase) RestoreProfile(ctx context.Context, userID uuid.UUID) error {
	if err := u.repo.UnarchiveProfile(ctx, userID); err != nil {
		return err
	}
	return nil
}
