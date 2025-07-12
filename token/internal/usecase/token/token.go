package token

import (
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	tokengen "mandacode.com/accounts/token/internal/infra/token"
)

type TokenUsecase struct {
	accessTokenGenerator            *tokengen.TokenGenerator
	refreshTokenGenerator           *tokengen.TokenGenerator
	emailVerificationTokenGenerator *tokengen.TokenGenerator
}

// GenerateAccessToken generates an access token for a user.
//
// Parameters:
//   - userID: The unique identifier of the user for whom the access token is generated.
//
// Returns:
//   - string: The generated JWT access token.
//   - int64: The expiration time of the token in seconds since epoch.
//   - error: An error if the token generation fails.
func (t *TokenUsecase) GenerateAccessToken(userID string) (string, int64, error) {
	claims := map[string]string{
		"sub": userID, // Use "sub" claim for user ID
	}
	return t.accessTokenGenerator.GenerateToken(claims)
}

// GenerateEmailVerificationToken generates an email verification token for a user.
//
// Parameters:
//   - userID: The unique identifier of the user for whom the email verification token is generated.
//   - email: The email address to be verified.
//   - code: The verification code to be included in the token.
//
// Returns:
//   - string: The generated JWT email verification token.
//   - int64: The expiration time of the token in seconds since epoch.
//   - error: An error if the token generation fails.
func (t *TokenUsecase) GenerateEmailVerificationToken(userID string, email string, code string) (string, int64, error) {
	claims := map[string]string{
		"sub":   userID,
		"email": email,
		"code":  code,
	}
	return t.emailVerificationTokenGenerator.GenerateToken(claims)
}

// GenerateRefreshToken generates a refresh token for a user.
//
// Parameters:
//   - userID: The unique identifier of the user for whom the refresh token is generated.
//
// Returns:
//   - string: The generated JWT refresh token.
//   - int64: The expiration time of the token in seconds since epoch.
//   - error: An error if the token generation fails.
func (t *TokenUsecase) GenerateRefreshToken(userID string) (string, int64, error) {
	claims := map[string]string{
		"sub": userID, // Use "sub" claim for user ID
	}
	return t.refreshTokenGenerator.GenerateToken(claims)
}

// VerifyAccessToken verifies the provided access token and returns the user ID if valid.
//
// Parameters:
//   - token: The JWT access token to be verified.
//
// Returns:
//   - *string: The user ID extracted from the token claims if verification is successful.
//   - error: An error if the token verification fails or if the user ID claim is missing.
func (t *TokenUsecase) VerifyAccessToken(token string) (*string, error) {
	claims, err := t.accessTokenGenerator.VerifyToken(token)
	if err != nil {
		joinedErr := errors.Join(err, "failed to verify access token")
		return nil, errors.Upgrade(joinedErr, errcode.ErrInvalidToken, "Token Verification Error")
	}

	userID, ok := claims["sub"]
	if !ok {
		return nil, errors.New("access token does not contain user ID claim", "Token Verification Error", errcode.ErrInvalidToken)
	}

	return &userID, nil
}

// VerifyEmailVerificationToken verifies the provided email verification token and returns the user ID, email, and code if valid.
// Parameters:
//   - token: The JWT email verification token to be verified.
//
// Returns:
//   - *string: The user ID extracted from the token claims if verification is successful.
//   - *string: The email extracted from the token claims if verification is successful.
//   - *string: The verification code extracted from the token claims if verification is successful.
//   - error: An error if the token verification fails or if any required claims are missing.
func (t *TokenUsecase) VerifyEmailVerificationToken(token string) (*string, *string, *string, error) {
	claims, err := t.emailVerificationTokenGenerator.VerifyToken(token)
	if err != nil {
		joinedErr := errors.Join(err, "failed to verify email verification token")
		return nil, nil, nil, errors.Upgrade(joinedErr, errcode.ErrInvalidToken, "Token Verification Error")
	}

	userID, ok := claims["sub"]
	if !ok {
		return nil, nil, nil, errors.New("email verification token does not contain user ID claim", "Token Verification Error", errcode.ErrInvalidToken)
	}

	email, ok := claims["email"]
	if !ok {
		return nil, nil, nil, errors.New("email verification token does not contain email claim", "Token Verification Error", errcode.ErrInvalidToken)
	}

	code, ok := claims["code"]
	if !ok {
		return nil, nil, nil, errors.New("email verification token does not contain code claim", "Token Verification Error", errcode.ErrInvalidToken)
	}

	return &userID, &email, &code, nil
}

// VerifyRefreshToken verifies the provided refresh token and returns the user ID if valid.
//
// Parameters:
//   - token: The JWT refresh token to be verified.
//
// Returns:
//   - *string: The user ID extracted from the token claims if verification is successful.
//   - error: An error if the token verification fails or if the user ID claim is missing.
func (t *TokenUsecase) VerifyRefreshToken(token string) (*string, error) {
	claims, err := t.refreshTokenGenerator.VerifyToken(token)
	if err != nil {
		joinedErr := errors.Join(err, "failed to verify refresh token")
		return nil, errors.Upgrade(joinedErr, errcode.ErrInvalidToken, "Token Verification Error")
	}

	userID, ok := claims["sub"]
	if !ok {
		return nil, errors.New("refresh token does not contain user ID claim", "Token Verification Error", errcode.ErrInvalidToken)
	}

	return &userID, nil
}

// NewTokenUsecase creates a new instance of tokenUsecase with the provided TokenGenerators.
func NewTokenUsecase(
	accessTokenGenerator *tokengen.TokenGenerator,
	refreshTokenGenerator *tokengen.TokenGenerator,
	emailVerificationTokenGenerator *tokengen.TokenGenerator,
) *TokenUsecase {
	return &TokenUsecase{
		accessTokenGenerator:            accessTokenGenerator,
		refreshTokenGenerator:           refreshTokenGenerator,
		emailVerificationTokenGenerator: emailVerificationTokenGenerator,
	}
}
