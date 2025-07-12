package tokenmodels

import "github.com/google/uuid"

type EmailVerificationResult struct {
	Valid  bool      `json:"valid"`
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Code   string    `json:"code"`
}
