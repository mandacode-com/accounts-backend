package dbrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/profile/ent"
	"mandacode.com/accounts/profile/ent/profile"
)

type ProfileRepository struct {
	client *ent.Client
}

func NewProfileRepository(client *ent.Client) *ProfileRepository {
	return &ProfileRepository{
		client: client,
	}
}

func (r *ProfileRepository) GetProfile(ctx context.Context, id uuid.UUID) (*ent.Profile, error) {
	prof, err := r.client.Profile.Query().Where(profile.UserID(id)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("profile not found", "Failed to get profile", errcode.ErrNotFound)
		}
		return nil, errors.Upgrade(err, "Failed to get profile", errcode.ErrInternalFailure)
	}
	return prof, nil
}

func (r *ProfileRepository) CreateProfile(ctx context.Context, data *CreateProfileModel) (*ent.Profile, error) {
	prof, err := r.client.Profile.Create().
		SetUserID(data.UserID).
		SetEmail(data.Email).
		SetNickname(data.Nickname).
		Save(ctx)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to create profile", errcode.ErrInternalFailure)
	}

	return prof, nil
}

func (r *ProfileRepository) UpdateProfile(ctx context.Context, data *UpdateProfileModel) (*ent.Profile, error) {
	prof, err := r.client.Profile.Query().Where(profile.UserID(data.UserID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("profile not found", "Failed to update profile", errcode.ErrNotFound)
		}
		return nil, errors.Upgrade(err, "Failed to update profile", errcode.ErrInternalFailure)
	}

	update := r.client.Profile.UpdateOne(prof)

	if data.Email != nil {
		update.SetEmail(*data.Email)
	}
	if data.Avatar != nil {
		update.SetAvatar(*data.Avatar)
	}
	if data.Bio != nil {
		update.SetBio(*data.Bio)
	}
	if data.Location != nil {
		update.SetLocation(*data.Location)
	}
	if data.Nickname != nil {
		update.SetNickname(*data.Nickname)
	}

	prof, err = update.Save(ctx)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to update profile", errcode.ErrInternalFailure)
	}

	return prof, nil
}

func (r *ProfileRepository) ArchiveProfile(ctx context.Context, userID uuid.UUID) error {
	prof, err := r.client.Profile.Query().Where(profile.UserID(userID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("profile not found", "Failed to archive profile", errcode.ErrNotFound)
		}
		return errors.Upgrade(err, "Failed to archive profile", errcode.ErrInternalFailure)
	}

	_, err = r.client.Profile.UpdateOne(prof).
		SetIsArchived(true).
		SetArchivedAt(prof.UpdatedAt).
		Save(ctx)
	if err != nil {
		return errors.Upgrade(err, "Failed to archive profile", errcode.ErrInternalFailure)
	}

	return nil
}

func (r *ProfileRepository) UnarchiveProfile(ctx context.Context, userID uuid.UUID) error {
	prof, err := r.client.Profile.Query().Where(profile.UserID(userID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("profile not found", "Failed to unarchive profile", errcode.ErrNotFound)
		}
		return errors.Upgrade(err, "Failed to unarchive profile", errcode.ErrInternalFailure)
	}

	_, err = r.client.Profile.UpdateOne(prof).
		SetIsArchived(false).
		SetNillableIsArchived(nil).
		Save(ctx)
	if err != nil {
		return errors.Upgrade(err, "Failed to unarchive profile", errcode.ErrInternalFailure)
	}

	return nil
}

func (r *ProfileRepository) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	prof, err := r.client.Profile.Query().Where(profile.UserID(userID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("profile not found", "Failed to delete profile", errcode.ErrNotFound)
		}
		return errors.Upgrade(err, "Failed to delete profile", errcode.ErrInternalFailure)
	}

	err = r.client.Profile.DeleteOne(prof).Exec(ctx)
	if err != nil {
		return errors.Upgrade(err, "Failed to delete profile", errcode.ErrInternalFailure)
	}

	return nil
}
