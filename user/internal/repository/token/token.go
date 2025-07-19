package tokenrepo

import (
	"context"

	"github.com/google/uuid"
	tokenv1 "github.com/mandacode-com/accounts-proto/go/token/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	tokenrepodto "mandacode.com/accounts/user/internal/repository/token/dto"
)

type MailTokenRepository struct {
	client tokenv1.TokenServiceClient
}

// GenerateEmailVerificationToken creates a new email verification token for the user.
//
// Parameters:
//   - ctx: The context for the operation.
//   - userID: The ID of the user for whom the email verification token is generated.
//   - email: The email address to verify.
//   - code: The verification code associated with the email.
func (t *MailTokenRepository) GenerateEmailVerificationToken(ctx context.Context, userID uuid.UUID, email string, code string) (string, int64, error) {
	resp, err := t.client.GenerateEmailVerificationToken(ctx, &tokenv1.GenerateEmailVerificationTokenRequest{
		UserId: userID.String(),
		Email:  email,
		Code:   code,
	})
	if err != nil {
		return "", 0, errors.Upgrade(err, "Failed to generate email verification token", errcode.ErrInternalFailure)
	}
	if err := resp.ValidateAll(); err != nil {
		return "", 0, errors.Upgrade(err, "Invalid response from token service", errcode.ErrInternalFailure)
	}
	return resp.Token, resp.ExpiresAt, nil
}

// VerifyEmailVerificationToken checks if the provided email verification token is valid.
//
// Parameters:
//   - ctx: The context for the operation.
//   - token: The email verification token to verify.
//
// Returns:
//   - data: A pointer to an EmailVerificationResult containing the verification result.
//   - error: An error if the verification fails, otherwise nil.
func (t *MailTokenRepository) VerifyEmailVerificationToken(ctx context.Context, token string) (*tokenrepodto.EmailVerificationResult, error) {
	resp, err := t.client.VerifyEmailVerificationToken(ctx, &tokenv1.VerifyEmailVerificationTokenRequest{Token: token})
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to verify email verification token", errcode.ErrInternalFailure)
	}
	if err := resp.ValidateAll(); err != nil {
		return nil, errors.Upgrade(err, "Invalid response from token service", errcode.ErrInternalFailure)
	}

	userUUID, err := uuid.Parse(*resp.UserId)
	if err != nil {
		return nil, errors.Upgrade(err, "Invalid user ID in response", errcode.ErrInternalFailure)
	}
	data := &tokenrepodto.EmailVerificationResult{
		Valid:  resp.Valid,
		UserID: userUUID,
		Email:  *resp.Email,
		Code:   *resp.Code,
	}
	return data, nil
}

func NewTokenRepository(client tokenv1.TokenServiceClient) *MailTokenRepository {
	return &MailTokenRepository{client: client}
}
