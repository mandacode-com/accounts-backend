package tokengen

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/token/internal/util"
)

// jwtGenerator is the concrete implementation of TokenGenerator
type TokenGenerator struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	expiresIn  time.Duration
}

// NewTokenGenerator creates a new tokenGenerator instance with the provided RSA keys and expiration duration
//
// Parameters:
//   - privateKey: the RSA private key used for signing the token
//   - expiresIn: the duration after which the token will expire
//
// Returns:
//   - svcdomain.TokenGenerator: an instance of TokenGenerator
//   - error: an error if the private key is nil or expiresIn is not greater than zero
func NewTokenGenerator(
	privateKey *rsa.PrivateKey,
	expiresIn time.Duration) (*TokenGenerator, error) {
	if privateKey == nil {
		return nil, errors.New("private key cannot be nil", "Invalid Private Key", errcode.ErrInvalidFormat)
	}
	if expiresIn <= 0 {
		return nil, errors.New("expiresIn must be greater than zero", "Invalid Expiration Duration", errcode.ErrInvalidFormat)
	}

	return &TokenGenerator{
		publicKey:  &privateKey.PublicKey,
		privateKey: privateKey,
		expiresIn:  expiresIn,
	}, nil
}

// NewTokenGeneratorByStr creates a new tokenGenerator using RSA keys provided as PEM formatted strings
//
// Parameters:
//   - privateKeyStr: the PEM formatted RSA private key string used for signing the token
//   - expiresIn: the duration after which the token will expiresIn
//
// Returns:
//   - svcdomain.TokenGenerator: an instance of TokenGenerator
//   - error: an error if the private key string is invalid or expiresIn is not greater than zero
func NewTokenGeneratorByStr(
	privateKeyStr string,
	expiresIn time.Duration) (*TokenGenerator, error) {
	privateKey, err := util.LoadRSAPrivateKeyFromPEM(privateKeyStr)

	if err != nil {
		return nil, err
	}
	if privateKey == nil {
		return nil, errors.New("private key cannot be nil", "Invalid Private Key", errcode.ErrInvalidFormat)
	}
	if expiresIn <= 0 {
		return nil, errors.New("expiresIn must be greater than zero", "Invalid Expiration Duration", errcode.ErrInvalidFormat)
	}

	return &TokenGenerator{
		publicKey:  &privateKey.PublicKey,
		privateKey: privateKey,
		expiresIn:  expiresIn,
	}, nil
}

func (j *TokenGenerator) GenerateToken(
	claims map[string]string,
) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(j.expiresIn)

	tokenClaims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": expiresAt.Unix(),
	}

	for key, value := range claims {
		tokenClaims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenClaims)
	signedToken, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", 0, errors.New(err.Error(), "Failed to sign token", errcode.ErrInternalFailure)
	}

	return signedToken, expiresAt.Unix(), nil
}

func (j *TokenGenerator) VerifyToken(
	token string,
) (map[string]string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method", "Invalid Token Signing Method", errcode.ErrInvalidToken)
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, errors.New(err.Error(), "Failed to parse token", errcode.ErrInvalidToken)
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		result := make(map[string]string)
		for key, value := range claims {
			if strValue, ok := value.(string); ok {
				result[key] = strValue
			}
		}
		return result, nil
	}

	return nil, errors.New("invalid token", "Token Verification Failed", errcode.ErrInvalidToken)
}
