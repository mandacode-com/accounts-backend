package signupinfra

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/go-playground/validator/v10"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/auth/ent/authaccount"
	signupinfradto "mandacode.com/accounts/auth/internal/infra/signup/dto"
)

type SignupAPI struct {
	endpoint  *url.URL
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
	urlEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Upgrade(err, "Failed to parse endpoint URL", errcode.ErrInvalidInput)
	}

	return &SignupAPI{
		endpoint:  urlEndpoint,
		client:    client,
		validator: validator,
	}, nil
}

func (s *SignupAPI) OAuthSignup(
	provider authaccount.Provider,
	accessToken string,
) (*signupinfradto.OAuthSignupResponse, error) {
	endpoint := s.endpoint
	endpoint.Path = path.Join(endpoint.Path, provider.String())
	req, err := http.NewRequest("GET", endpoint.String(), nil)
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
