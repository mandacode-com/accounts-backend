package token

import (
	"context"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	tokenrepo "mandacode.com/accounts/auth/internal/repository/token"
)

type RefreshUsecase struct {
	token *tokenrepo.TokenRepository
}

// Refresh generates new access and refresh tokens based on a valid refresh token.
//
// Parameters:
//   - ctx: The context for the operation.
//
// Returns:
//   - newAccessToken: The newly generated access token.
//   - newRefreshToken: The newly generated refresh token.
//   - err: An error if the operation fails, or nil if successful.
func (r *RefreshUsecase) Refresh(ctx context.Context, refreshToken string) (newAccessToken string, newRefreshToken string, err error) {
	// Validate the refresh token
	valid, userID, err := r.token.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", errors.New("failed to verify refresh token", "Unauthorized", errcode.ErrUnauthorized)
	}
	if !valid || userID == nil {
		return "", "", errors.New("invalid refresh token", "Unauthorized", errcode.ErrUnauthorized)
	}

	// Generate new access and refresh tokens
	userUID, err := uuid.Parse(*userID)
	if err != nil {
		return "", "", errors.New("invalid user ID in refresh token", "Invalid User ID", errcode.ErrUnauthorized)
	}
	newAccessToken, _, err = r.token.GenerateAccessToken(ctx, userUID)
	if err != nil {
		return "", "", errors.Join(err, "failed to generate new access token")
	}
	newRefreshToken, _, err = r.token.GenerateRefreshToken(ctx, userUID)
	if err != nil {
		return "", "", errors.Join(err, "failed to generate new refresh token")
	}

	return newAccessToken, newRefreshToken, nil
}

// NewRefreshUsecase creates a new instance of RefreshUsecase with the provided token repository.
func NewRefreshUsecase(token *tokenrepo.TokenRepository) *RefreshUsecase {
	return &RefreshUsecase{
		token: token,
	}
}
