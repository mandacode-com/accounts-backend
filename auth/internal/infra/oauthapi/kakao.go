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

type KakaoAPI struct {
	clientID     string
	clientSecret string
	redirectURL  string
	validator    *validator.Validate
}

// GetUserInfo fetches user information from Kakao using the provided access token.
func (k *KakaoAPI) GetUserInfo(accessToken string) (*oauthmodels.UserInfo, error) {
	req, err := http.NewRequest("GET", oauthapimeta.KakaoUserInfoEndpoint, nil)
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
	var rawUserInfo infoapidto.RawKakaoUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&rawUserInfo); err != nil {
		return nil, errors.New("failed to decode user info: " + err.Error())
	}

	oauthUserInfo := oauthmodels.NewUserInfo(
		rawUserInfo.ID,
		rawUserInfo.Email,
		rawUserInfo.Name,
		rawUserInfo.EmailVerified,
	)
	if err := k.validator.Struct(oauthUserInfo); err != nil {
		return nil, errors.New("invalid user info structure: " + err.Error())
	}

	return oauthUserInfo, nil
}

// NewKakaoAPI creates a new instance of KakaoAPI with the required parameters.
func NewKakaoAPI(clientID, clientSecret, redirectURL string, validator *validator.Validate) (OAuthAPI, error) {
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return nil, errors.New("clientID, clientSecret, and redirectURL must not be empty")
	}

	return &KakaoAPI{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		validator:    validator,
	}, nil
}

func (k *KakaoAPI) GetAccessToken(code string) (string, error) {
	req, err := http.NewRequest("POST", oauthapimeta.KakaoTokenEndpoint, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("code", code)
	q.Add("client_id", k.clientID)
	q.Add("client_secret", k.clientSecret)
	q.Add("redirect_uri", k.redirectURL)
	q.Add("grant_type", oauthapimeta.KakaoGrantType)
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

	var tokenResponse codeapidto.KakaoAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", errors.New("failed to decode access token response: " + err.Error())
	}
	if tokenResponse.AccessToken == "" {
		return "", errors.New("access token is empty")
	}

	return tokenResponse.AccessToken, nil
}

func (k *KakaoAPI) GetLoginURL() string {
	loginURL := oauthapimeta.KakaoAuthEndpoint + "?client_id=" + k.clientID +
		"&redirect_uri=" + k.redirectURL + "&response_type=code&scope=account_email%20profile_nickname"

	return loginURL
}
