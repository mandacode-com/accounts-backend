package signup

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	authrepo "mandacode.com/accounts/user/internal/repository/auth"
	authrepodto "mandacode.com/accounts/user/internal/repository/auth/dto"
	coderepo "mandacode.com/accounts/user/internal/repository/code"
	dbrepo "mandacode.com/accounts/user/internal/repository/database"
	maileventrepo "mandacode.com/accounts/user/internal/repository/mailevent"
	tokenrepo "mandacode.com/accounts/user/internal/repository/token"
)

type VerifyEmailUsecase struct {
	sentEmailRepo         *dbrepo.SentEmailRepository
	authRepo              *authrepo.AuthRepository
	mailTokenRepo         *tokenrepo.MailTokenRepository
	mailEventEmitter      *maileventrepo.MailEventEmitter
	codeManager           *coderepo.CodeManager
	emailVerificationLink string
	maxSentEmails         int
	maxSentEmailsDuration time.Duration
}

// NewVerifyEmailUsecase creates a new instance of VerifyEmailUsecase with the provided repositories.
func NewVerifyEmailUsecase(
	sentEmailRepo *dbrepo.SentEmailRepository,
	authRepo *authrepo.AuthRepository,
	mailTokenRepo *tokenrepo.MailTokenRepository,
	mailEventEmitter *maileventrepo.MailEventEmitter,
	codeManager *coderepo.CodeManager,
	emailVerificationLink string,
	maxSentEmails int,
	maxSentEmailsDuration time.Duration,
) *VerifyEmailUsecase {
	return &VerifyEmailUsecase{
		sentEmailRepo:         sentEmailRepo,
		authRepo:              authRepo,
		mailTokenRepo:         mailTokenRepo,
		mailEventEmitter:      mailEventEmitter,
		codeManager:           codeManager,
		emailVerificationLink: emailVerificationLink,
		maxSentEmails:         maxSentEmails,
		maxSentEmailsDuration: maxSentEmailsDuration,
	}
}

// canSendVerificationEmail checks if the user can send a verification email.
func (v *VerifyEmailUsecase) canSendVerificationEmail(ctx context.Context, userID uuid.UUID) (bool, error) {
	count, err := v.sentEmailRepo.GetSentEmailNumberByUserDuration(ctx, userID, v.maxSentEmailsDuration)
	if err != nil {
		return false, errors.Upgrade(err, "Failed to get sent emails by user ID", errcode.ErrInternalFailure)
	}

	if count >= v.maxSentEmails {
		return false, nil
	}
	return true, nil
}

// SendVerificationEmail sends a verification email to the user.
func (v *VerifyEmailUsecase) SendVerificationEmail(ctx context.Context, userID uuid.UUID, email string) error {
	// Check if Verification Email can be sent
	canSend, err := v.canSendVerificationEmail(ctx, userID)
	if err != nil {
		return errors.Upgrade(err, "Failed to check if verification email can be sent", errcode.ErrInternalFailure)
	}
	if !canSend {
		return errors.New("Too many verification emails sent", "You have reached the maximum number of verification emails sent", errcode.ErrTooManyRequests)
	}

	// Issue a verification code and generate a token
	code, err := v.codeManager.IssueCode(ctx, userID)
	if err != nil {
		return errors.Upgrade(err, "Failed to issue verification code", errcode.ErrInternalFailure)
	}
	token, _, err := v.mailTokenRepo.GenerateEmailVerificationToken(ctx, userID, email, code)
	if err != nil {
		return errors.Upgrade(err, "Failed to generate email verification token", errcode.ErrInternalFailure)
	}
	verificationLink := v.emailVerificationLink + "?token=" + token
	if err := v.mailEventEmitter.SendEmailVerificationMail(email, verificationLink); err != nil {
		return errors.Upgrade(err, "Failed to send email verification mail", errcode.ErrInternalFailure)
	}
	_, err = v.sentEmailRepo.CreateSentEmail(ctx, userID, email)
	if err != nil {
		return errors.Upgrade(err, "Failed to create sent email record", errcode.ErrInternalFailure)
	}
	return nil
}

// VerifyEmailToken verifies the email verification token and returns the user ID if valid.
func (v *VerifyEmailUsecase) VerifyEmailToken(ctx context.Context, token string) (uuid.UUID, error) {
	result, err := v.mailTokenRepo.VerifyEmailVerificationToken(ctx, token)
	if err != nil {
		return uuid.Nil, errors.Upgrade(err, "Failed to verify email verification token", errcode.ErrInternalFailure)
	}
	if result == nil {
		return uuid.Nil, errors.New("Invalid email verification token", "The provided email verification token is invalid", errcode.ErrInvalidToken)
	}
	userID, err := uuid.Parse(result.UserID.String())
	if err != nil {
		return uuid.Nil, errors.Upgrade(err, "Invalid user ID in response", errcode.ErrInternalFailure)
	}
	valid, err := v.codeManager.ValidateCode(ctx, userID, result.Code)
	if err != nil {
		return uuid.Nil, errors.Upgrade(err, "Failed to validate verification code", errcode.ErrInternalFailure)
	}
	if !valid {
		return uuid.Nil, errors.New("Invalid verification code", "The provided verification code is invalid", errcode.ErrInvalidToken)
	}

	_, err = v.authRepo.UpdateLocalUserEmailVerification(ctx, &authrepodto.UpdateEmailVerificationRequest{
		UserID:   userID,
		Verified: true,
	})
	if err != nil {
		return uuid.Nil, errors.Upgrade(err, "Failed to update user email verification status", errcode.ErrInternalFailure)
	}

	return userID, nil
}
