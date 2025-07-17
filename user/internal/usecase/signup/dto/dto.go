package signupdto

import (
	"time"

	"github.com/google/uuid"
	"mandacode.com/accounts/user/internal/models/provider"
)

type LocalSignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LocalSignupResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type OAuthSignupRequest struct {
	Provider    provider.ProviderType `json:"provider"`
	AccessToken string                `json:"access_token,omitempty"`
}

type OAuthSignupResponse struct {
	UserID    uuid.UUID             `json:"user_id"`
	Provider  provider.ProviderType `json:"provider"`
	Email     string                `json:"email"`
	CreatedAt time.Time             `json:"created_at"`
}

type SendVerificationEmailRequest struct {
	Email string `json:"email"`
}
