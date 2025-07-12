package dto

import (
	"github.com/google/uuid"
	dbrepo "mandacode.com/accounts/profile/internal/repository/database"
)

type UpdateProfileData struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    *string   `json:"email,omitempty"`
	Avatar   *string   `json:"avatar,omitempty"`
	Bio      *string   `json:"bio,omitempty"`
	Location *string   `json:"location,omitempty"`
	Nickname *string   `json:"nickname,omitempty"`
}

func (data UpdateProfileData) ToRepoModel() *dbrepo.UpdateProfileModel {
	return &dbrepo.UpdateProfileModel{
		UserID:   data.UserID,
		Email:    data.Email,
		Avatar:   data.Avatar,
		Bio:      data.Bio,
		Location: data.Location,
		Nickname: data.Nickname,
	}
}
