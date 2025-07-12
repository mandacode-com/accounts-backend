package usermodels

import (
	"time"

	"github.com/google/uuid"
	"mandacode.com/accounts/user/ent"
)

type SecureUser struct {
	ID          uuid.UUID  `json:"id"`
	SyncCode    string     `json:"sync_code"`
	IsActive    bool       `json:"is_active"`
	IsBlocked   bool       `json:"is_blocked"`
	IsArchived  bool       `json:"is_archived"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeleteAfter *time.Time `json:"delete_after,omitempty"`
}

// NewSecureUser creates a new SecureUser with the current time for CreatedAt and UpdatedAt.
func NewSecureUser(user *ent.User) *SecureUser {
	return &SecureUser{
		ID:        user.ID,
		SyncCode:  user.SyncCode,
		IsActive:  user.IsActive,
		IsBlocked: user.IsBlocked,
		IsArchived: user.IsArchived,
		ArchivedAt:  user.ArchivedAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeleteAfter: user.DeleteAfter,
	}
}
