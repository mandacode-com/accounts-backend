package oauthapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	codeapidto "mandacode.com/accounts/auth/internal/infra/oauthapi/dto/codeapi"
	infoapidto "mandacode.com/accounts/auth/internal/infra/oauthapi/dto/infoapi"
	oauthapimeta "mandacode.com/accounts/auth/internal/infra/oauthapi/meta"
	oauthmodels "mandacode.com/accounts/auth/internal/models/oauth"
)

type googleAPI struct {
	clientID     string
	clientSecret string
	redirectURL  string
	validator    *validator.Validate
}

// GetUserInfo fetches user information from Google using the provided access token.
func (g *googleAPI) GetUserInfo(accessToken string) (*oauthmodels.UserInfo, error) {
	req, err := http.NewRequest("GET", oauthapimeta.GoogleUserInfoEndpoint, nil)
	if err != nil {
		return nil, errors.New("failed to create request: " + err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("failed to fetch user info: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch user info: status code " + resp.Status)
	}

	// Decode the JSON response into the RawGoogleUserInfo struct
	var rawUserInfo infoapidto.RawGoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&rawUserInfo); err != nil {
		return nil, errors.New("failed to decode user info: " + err.Error())
	}

	oauthUserInfo := oauthmodels.NewUserInfo(
		rawUserInfo.Sub,
		rawUserInfo.Email,
		rawUserInfo.Name,
		rawUserInfo.EmailVerified,
	)
	if err := g.validator.Struct(oauthUserInfo); err != nil {
		return nil, errors.New("invalid user info structure: " + err.Error())
	}

	return oauthUserInfo, nil
}

// NewGoogleAPI creates a new instance of GoogleAPI with the required parameters.
func NewGoogleAPI(clientID, clientSecret, redirectURL string, validator *validator.Validate) (OAuthAPI, error) {
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return nil, errors.New("client ID, client secret, and redirect URL must be set")
	}

	return &googleAPI{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		validator:    validator,
	}, nil
}

func (g *googleAPI) GetAccessToken(code string) (string, error) {
	req, err := http.NewRequest("POST", oauthapimeta.GoogleTokenEndpoint, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("code", code)
	q.Add("client_id", g.clientID)
	q.Add("client_secret", g.clientSecret)
	q.Add("redirect_uri", g.redirectURL)
	q.Add("grant_type", oauthapimeta.GoogleGrantType)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to retrieve access token: " + resp.Status)
	}

	var tokenResponse codeapidto.GoogleAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", errors.New("failed to decode access token response: " + err.Error())
	}
	if tokenResponse.AccessToken == "" {
		return "", errors.New("access token is empty in response")
	}

	return tokenResponse.AccessToken, nil
}

func (g *googleAPI) GetLoginURL() string {
	loginUrl := oauthapimeta.GoogleAuthEndpoint + "?client_id=" + g.clientID +
		"&redirect_uri=" + g.redirectURL +
		"&response_type=code" +
		"&scope=email%20profile" +
		"&access_type=offline"

	return loginUrl
}
