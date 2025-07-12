package dbmodels

import (
	"time"

	"github.com/google/uuid"
	"mandacode.com/accounts/profile/ent"
)

type SecureProfile struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewSecureProfile(entProfile *ent.Profile) *SecureProfile {
	if entProfile == nil {
		return nil
	}
	return &SecureProfile{
		UserID:    entProfile.UserID,
		Email:     entProfile.Email,
		Avatar:    entProfile.Avatar,
		Bio:       entProfile.Bio,
		CreatedAt: entProfile.CreatedAt,
		UpdatedAt: entProfile.UpdatedAt,
	}
}
