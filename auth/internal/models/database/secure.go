package dbmodels

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/authaccount"
)

type SecureLocalAuthAccount struct {
	ID         uuid.UUID            `json:"id"`
	UserID     uuid.UUID            `json:"user_id" validate:"required"`
	Provider   authaccount.Provider `json:"provider" validate:"required,oneof=local"`
	Email      string               `json:"email" validate:"required,email"`
	IsVerified bool                 `json:"is_verified" validate:"required"`
}

func NewSecureLocalAuthAccount(authAccount *ent.AuthAccount) *SecureLocalAuthAccount {
	return &SecureLocalAuthAccount{
		ID:         authAccount.ID,
		UserID:     authAccount.UserID,
		Provider:   authAccount.Provider,
		Email:      authAccount.Email,
		IsVerified: authAccount.IsVerified,
	}
}

type SecureOAuthAuthAccount struct {
	ID         uuid.UUID            `json:"id"`
	UserID     uuid.UUID            `json:"user_id" validate:"required"`
	Provider   authaccount.Provider `json:"provider" validate:"required,oneof=google kakao naver apple"`
	ProviderID string               `json:"provider_id" validate:"required"`
	Email      string               `json:"email" validate:"required,email"`
	IsVerified bool                 `json:"is_verified" validate:"required"`
}

func NewSecureOAuthAuthAccount(authAccount *ent.AuthAccount) *SecureOAuthAuthAccount {
	return &SecureOAuthAuthAccount{
		ID:         authAccount.ID,
		UserID:     authAccount.UserID,
		Provider:   authAccount.Provider,
		ProviderID: *authAccount.ProviderID,
		Email:      authAccount.Email,
		IsVerified: authAccount.IsVerified,
	}
}

type SecureAuthAccount struct {
	ID         uuid.UUID            `json:"id"`
	UserID     uuid.UUID            `json:"user_id" validate:"required"`
	Provider   authaccount.Provider `json:"provider" validate:"required,oneof=local google kakao naver apple"`
	Email      string               `json:"email" validate:"required,email"`
	IsVerified bool                 `json:"is_verified" validate:"required"`
}

func NewSecureAuthAccount(authAccount *ent.AuthAccount) *SecureAuthAccount {
	return &SecureAuthAccount{
		ID:         authAccount.ID,
		UserID:     authAccount.UserID,
		Provider:   authAccount.Provider,
		Email:      authAccount.Email,
		IsVerified: authAccount.IsVerified,
	}
}
