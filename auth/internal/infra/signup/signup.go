package signupinfra

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	signupinfradto "mandacode.com/accounts/auth/internal/infra/signup/dto"
)

type SignupAPI struct {
	endpoint  string
	client    *http.Client
	validator *validator.Validate
}

// NewSignup creates a new Signup instance with the provided HTTP client.
func NewSignupApi(endpoint string, client *http.Client, validator *validator.Validate) (*SignupAPI, error) {
	if endpoint == "" {
		return nil, errors.New("endpoint cannot be empty", "InvalidEndpoint", errcode.ErrInvalidInput)
	}
	if client == nil {
		return nil, errors.New("HTTP client cannot be nil", "InvalidClient", errcode.ErrInvalidInput)
	}
	if validator == nil {
		return nil, errors.New("validator cannot be nil", "InvalidValidator", errcode.ErrInvalidInput)
	}

	return &SignupAPI{
		endpoint:  endpoint,
		client:    client,
		validator: validator,
	}, nil
}

func (s *SignupAPI) OAuthSignup(
	accessToken string,
) (*signupinfradto.OAuthSignupResponse, error) {
	req, err := http.NewRequest("POST", s.endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to sign up with OAuth provider", "OAuth Signup Error", errcode.ErrInternalFailure)
	}

	var rawResponse signupinfradto.OAuthSignupResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, errors.Upgrade(err, "Failed to decode OAuth signup response", errcode.ErrInternalFailure)
	}
	if err := s.validator.Struct(rawResponse); err != nil {
		return nil, errors.Upgrade(err, "Invalid OAuth signup response structure", errcode.ErrInvalidInput)
	}

	return &rawResponse, nil
}
