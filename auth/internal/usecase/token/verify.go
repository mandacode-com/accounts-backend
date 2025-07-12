package token

import (
	"context"

	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	tokenrepo "mandacode.com/accounts/auth/internal/repository/token"
)

type VerifyUsecase struct {
	token *tokenrepo.TokenRepository
}

// Verify verifies the access token and returns whether it is valid, the user ID if valid, or an error if verification fails.
//
// Parameters:
//   - ctx: The context for the operation.
//   - token: The access token to be verified.
//
// Returns:
//   - valid: A boolean indicating whether the token is valid.
//   - userID: The user ID associated with the token if valid, or nil if invalid.
//   - err: An error if the verification fails, or nil if successful.
func (v *VerifyUsecase) Verify(ctx context.Context, token string) (valid bool, userID *string, err error) {
	valid, userID, err = v.token.VerifyAccessToken(ctx, token)
	if err != nil {
		joinedErr := errors.Join(err, "failed to verify access token")
		return false, nil, errors.Upgrade(joinedErr, "Unauthorized", errcode.ErrUnauthorized)
	}
	if !valid || userID == nil {
		return false, nil, nil // Token is invalid or user ID is not present
	}
	return true, userID, nil
}

// VerifyRefresh verifies the refresh token and returns whether it is valid, the user ID if valid, or an error if verification fails.
//
// Parameters:
//   - ctx: The context for the operation.
//   - token: The refresh token to be verified.
//
// Returns:
//   - valid: A boolean indicating whether the token is valid.
//   - userID: The user ID associated with the token if valid, or nil if invalid.
func (v *VerifyUsecase) VerifyRefresh(ctx context.Context, token string) (valid bool, userID *string, err error) {
	valid, userID, err = v.token.VerifyRefreshToken(ctx, token)
	if err != nil {
		joinedErr := errors.Join(err, "failed to verify refresh token")
		return false, nil, errors.Upgrade(joinedErr, "Unauthorized", errcode.ErrUnauthorized)
	}
	if !valid || userID == nil {
		return false, nil, nil // Token is invalid or user ID is not present
	}
	return true, userID, nil
}

// NewVerifyUsecase creates a new instance of VerifyUsecase.
func NewVerifyUsecase(token *tokenrepo.TokenRepository) *VerifyUsecase {
	return &VerifyUsecase{
		token: token,
	}
}
