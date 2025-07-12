package dbmodels

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent/authaccount"
)

type CreateLocalAuthAccountInput struct {
	UserID     uuid.UUID `json:"user_id" validate:"required"`
	Email      string    `json:"email" validate:"required,email"`
	Password   string    `json:"password" validate:"required"`
	IsVerified bool      `json:"is_verified" validate:"omitempty"`
}

type CreateOAuthAuthAccountInput struct {
	UserID     uuid.UUID            `json:"user_id" validate:"required"`
	Provider   authaccount.Provider `json:"provider" validate:"required,oneof=local google kakao naver apple"`
	ProviderID string               `json:"provider_id" validate:"required"`
	Email      string               `json:"email" validate:"required,email"`
	IsVerified bool                 `json:"is_verified" validate:"omitempty"`
}
