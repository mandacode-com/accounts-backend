package localauth

import (
	"context"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/auth/internal/infra/mailer"
	dbmodels "mandacode.com/accounts/auth/internal/models/database"
	coderepo "mandacode.com/accounts/auth/internal/repository/code"
	dbrepo "mandacode.com/accounts/auth/internal/repository/database"
	tokenrepo "mandacode.com/accounts/auth/internal/repository/token"
	userrepo "mandacode.com/accounts/auth/internal/repository/user"
	localauthdto "mandacode.com/accounts/auth/internal/usecase/localauth/dto"
)

type SignupUsecase struct {
	authAccount      *dbrepo.AuthAccountRepository
	userService      *userrepo.UserServiceRepository
	token            *tokenrepo.TokenRepository
	mailer           *mailer.Mailer
	emailCodeManager *coderepo.CodeManager
	verifyEmailURL   string
}

// ResendVerificationEmail implements localauthdomain.SignupUsecase.
func (s *SignupUsecase) ResendVerificationEmail(ctx context.Context, email string) (success bool, err error) {
	auth, err := s.authAccount.GetLocalAuthAccountByEmail(ctx, email)
	if err != nil {
		return false, errors.Upgrade(err, "Unauthorized", errcode.ErrUnauthorized)
	}

	if auth.IsVerified {
		return false, errors.New("user already verified", "User Already Verified", errcode.ErrConflict)
	}

	// Generate email verification token
	code, err := s.emailCodeManager.IssueCode(ctx, auth.UserID)
	if err != nil {
		return false, errors.Upgrade(err, "Internal Error", errcode.ErrInternalFailure)
	}

	token, _, err := s.token.GenerateEmailVerificationToken(ctx, auth.UserID, email, code)
	if err != nil {
		return false, errors.Upgrade(err, "Internal Error", errcode.ErrInternalFailure)
	}

	// Send verification email
	url := s.verifyEmailURL + "?token=" + token
	err = s.mailer.SendEmailVerificationMail(email, url)
	if err != nil {
		return false, errors.Upgrade(err, "Failed to send verification email", errcode.ErrInternalFailure)
	}

	return true, nil
}

// Signup implements localauthdomain.SignupUsecase.
func (s *SignupUsecase) Signup(ctx context.Context, input localauthdto.SignupInput) (userID uuid.UUID, err error) {
	userID = uuid.New()
	createUserResp, err := s.userService.InitUser(ctx, userID)
	if err != nil {
		return uuid.Nil, errors.Upgrade(err, "Failed to create user", errcode.ErrInternalFailure)
	}
	if createUserResp == nil || createUserResp.UserId != userID.String() {
		s.userService.DeleteUser(ctx, userID)
		return uuid.Nil, errors.New("failed to create user", "Internal Error", errcode.ErrInternalFailure)
	}

	auth, err := s.authAccount.CreateLocalAuthAccount(
		ctx,
		&dbmodels.CreateLocalAuthAccountInput{
			UserID:     userID,
			Email:      input.Email,
			Password:   input.Password,
			IsVerified: false,
		},
	)
	if err != nil {
		s.userService.DeleteUser(ctx, userID)
		return uuid.Nil, errors.Join(err, "failed to create user")
	}

	// Generate email verification token
	code, err := s.emailCodeManager.IssueCode(ctx, auth.UserID)

	token, _, err := s.token.GenerateEmailVerificationToken(ctx, auth.UserID, input.Email, code)
	if err != nil {
		return uuid.Nil, errors.Upgrade(err, "Internal Error", errcode.ErrInternalFailure)
	}

	// Send verification email
	url := s.verifyEmailURL + "?token=" + token
	err = s.mailer.SendEmailVerificationMail(input.Email, url)
	if err != nil {
		return uuid.Nil, errors.Upgrade(err, "Failed to send verification email", errcode.ErrInternalFailure)
	}

	return auth.UserID, nil
}

// VerifyEmail implements localauthdomain.SignupUsecase.
func (s *SignupUsecase) VerifyEmail(ctx context.Context, email string, token string) (success bool, err error) {
	result, err := s.token.VerifyEmailVerificationToken(ctx, token)
	if err != nil {
		return false, errors.Upgrade(err, "Unauthorized", errcode.ErrUnauthorized)
	}
	if !result.Valid {
		return false, errors.New("invalid or expired token", "Unauthorized", errcode.ErrUnauthorized)
	}

	auth, err := s.authAccount.GetLocalAuthAccountByUserID(ctx, result.UserID)
	if err != nil {
		return false, errors.Upgrade(err, "Unauthorized", errcode.ErrUnauthorized)
	}

	if auth.Email != email {
		return false, errors.New("email does not match", "Unauthorized", errcode.ErrUnauthorized)
	}
	if auth.IsVerified {
		return false, errors.New("user already verified", "User Already Verified", errcode.ErrConflict)
	}
	// Verify the code
	valid, err := s.emailCodeManager.ValidateCode(ctx, auth.UserID, result.Code)
	if err != nil {
		return false, errors.Upgrade(err, "Failed to validate verification code", errcode.ErrInternalFailure)
	}
	if !valid {
		return false, errors.New("verification code is invalid or expired", "Unauthorized", errcode.ErrUnauthorized)
	}

	// Update the user's verification status
	_, err = s.authAccount.SetIsVerifiedByID(ctx, auth.ID, true)
	if err != nil {
		return false, errors.Upgrade(err, "Failed to update user verification status", errcode.ErrInternalFailure)
	}

	return true, nil
}

// NewSignupUsecase creates a new instance of SignupUsecase.
func NewSignupUsecase(
	authAccount *dbrepo.AuthAccountRepository,
	token *tokenrepo.TokenRepository,
	mailer *mailer.Mailer,
	emailCodeManager *coderepo.CodeManager,
	verifyEmailURL string,
) *SignupUsecase {
	return &SignupUsecase{
		authAccount:      authAccount,
		token:            token,
		mailer:           mailer,
		emailCodeManager: emailCodeManager,
		verifyEmailURL:   verifyEmailURL,
	}
}
