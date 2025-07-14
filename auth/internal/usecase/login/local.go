package login

import (
	"context"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"

	coderepo "mandacode.com/accounts/auth/internal/repository/code"
	dbrepo "mandacode.com/accounts/auth/internal/repository/database"
	tokenrepo "mandacode.com/accounts/auth/internal/repository/token"
	logindto "mandacode.com/accounts/auth/internal/usecase/login/dto"
)

type LocalLoginUsecase struct {
	authAccount      *dbrepo.AuthAccountRepository
	token            *tokenrepo.TokenRepository
	loginCodeManager *coderepo.CodeManager
}

func (l *LocalLoginUsecase) checkUserVerified(ctx context.Context, input logindto.LocalLoginInput) (uuid.UUID, error) {
	verified, userID, err := l.authAccount.ComparePassword(ctx, input.Email, input.Password)
	if err != nil {
		return uuid.Nil, err
	}
	if !verified {
		return uuid.Nil, errors.New("invalid email or password", "Unauthorized", errcode.ErrUnauthorized)
	}

	authAccount, err := l.authAccount.GetLocalAuthAccountByUserID(ctx, userID)
	if err != nil {
		return uuid.Nil, errors.Upgrade(err, "Failed to get auth account", errcode.ErrInternalFailure)
	}
	if !authAccount.IsVerified {
		return uuid.Nil, errors.New("user is not verified", "User Email Not Verified", errcode.ErrUnauthorized)
	}

	return userID, nil
}

// IssueLoginCode implements localauthdomain.LocalLoginUsecase.
func (l *LocalLoginUsecase) IssueLoginCode(ctx context.Context, input logindto.LocalLoginInput) (code string, userID uuid.UUID, err error) {
	userID, err = l.checkUserVerified(ctx, input)
	if err != nil {
		return "", uuid.Nil, err
	}

	code, err = l.loginCodeManager.IssueCode(ctx, userID)
	if err != nil {
		return "", uuid.Nil, errors.Upgrade(err, "Internal Error", errcode.ErrInternalFailure)
	}

	userID = userID
	return code, userID, nil
}

// VerifyLoginCode implements localauthdomain.LocalLoginUsecase.
func (l *LocalLoginUsecase) VerifyLoginCode(ctx context.Context, userID uuid.UUID, code string) (accessToken string, refreshToken string, err error) {
	valid, err := l.loginCodeManager.ValidateCode(ctx, userID, code)
	if err != nil {
		return "", "", errors.Upgrade(err, "Failed to validate login code", errcode.ErrInternalFailure)
	}
	if !valid {
		return "", "", errors.New("login code is invalid or expired", "Failed to validate login code", errcode.ErrUnauthorized)
	}

	// Generate access and refresh tokens
	return l.issueToken(ctx, userID)
}

// Login implements localauthdomain.LocalLoginUsecase.
func (l *LocalLoginUsecase) Login(ctx context.Context, input logindto.LocalLoginInput) (accessToken string, refreshToken string, err error) {
	userID, err := l.checkUserVerified(ctx, input)
	if err != nil {
		return "", "", err
	}

	// Generate access and refresh tokens
	return l.issueToken(ctx, userID)
}

// issueToken issues a new access token and refresh token for the user.
func (l *LocalLoginUsecase) issueToken(ctx context.Context, userID uuid.UUID) (accessToken string, refreshToken string, err error) {
	accessToken, _, err = l.token.GenerateAccessToken(ctx, userID)
	if err != nil {
		return "", "", errors.Upgrade(err, "Failed to generate token", errcode.ErrInternalFailure)
	}
	refreshToken, _, err = l.token.GenerateRefreshToken(ctx, userID)
	if err != nil {
		return "", "", errors.Upgrade(err, "Failed to generate token", errcode.ErrInternalFailure)
	}
	return accessToken, refreshToken, nil
}

func NewLocalLoginUsecase(
	authAccount *dbrepo.AuthAccountRepository,
	token *tokenrepo.TokenRepository,
	loginCodeManager *coderepo.CodeManager,
) *LocalLoginUsecase {
	return &LocalLoginUsecase{
		authAccount:      authAccount,
		token:            token,
		loginCodeManager: loginCodeManager,
	}
}
