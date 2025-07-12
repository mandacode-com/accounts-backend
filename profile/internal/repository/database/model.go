package dbrepo

import "github.com/google/uuid"

type CreateProfileModel struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email,omitempty"`
	Nickname string    `json:"nickname,omitempty"`
}

type UpdateProfileModel struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    *string   `json:"email,omitempty"`
	Avatar   *string   `json:"avatar,omitempty"`
	Bio      *string   `json:"bio,omitempty"`
	Location *string   `json:"location,omitempty"`
	Nickname *string   `json:"nickname,omitempty"`
}
