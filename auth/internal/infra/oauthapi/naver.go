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

type NaverAPI struct {
	clientID         string
	clientSecret     string
	redirectURL      string
	tokenRedirectURL string
	validator        *validator.Validate
}

// GetUserInfo implements oauthapidomain.OAuthCode.
func (n *NaverAPI) GetUserInfo(accessToken string) (*oauthmodels.UserInfo, error) {
	req, err := http.NewRequest("GET", oauthapimeta.NaverUserInfoEndpoint, nil)
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
	var rawUserInfo infoapidto.RawNaverUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&rawUserInfo); err != nil {
		return nil, errors.New("failed to decode user info: " + err.Error())
	}

	oauthUserInfo := oauthmodels.NewUserInfo(
		rawUserInfo.ID,
		rawUserInfo.Email,
		rawUserInfo.Name,
		true, // Naver does not provide email verification status
	)
	if err := n.validator.Struct(oauthUserInfo); err != nil {
		return nil, errors.New("invalid user info structure: " + err.Error())
	}

	return oauthUserInfo, nil
}

// NewNaverAPI creates a new instance of NaverAPI with the required parameters.
func NewNaverAPI(clientID, clientSecret, redirectURL string, validator *validator.Validate) (OAuthAPI, error) {
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return nil, errors.New("client ID, client secret, and redirect URL must be set")
	}

	return &NaverAPI{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		validator:    validator,
	}, nil
}

func (n *NaverAPI) GetAccessToken(code string) (string, error) {
	req, err := http.NewRequest("POST", oauthapimeta.NaverTokenEndpoint, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("code", code)
	q.Add("client_id", n.clientID)
	q.Add("client_secret", n.clientSecret)
	q.Add("redirect_uri", n.redirectURL)
	q.Add("grant_type", oauthapimeta.NaverGrantType)
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

	var tokenResponse codeapidto.NaverAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", errors.New("failed to decode access token response: " + err.Error())
	}
	if tokenResponse.AccessToken == "" {
		return "", errors.New("access token is empty")
	}

	return tokenResponse.AccessToken, nil
}

func (n *NaverAPI) GetLoginURL() string {
	loginURL := oauthapimeta.NaverAuthEndpoint + "?client_id=" + n.clientID +
		"&response_type=code" +
		"&redirect_uri=" + n.redirectURL +
		"&state=" + n.tokenRedirectURL

	return loginURL
}
