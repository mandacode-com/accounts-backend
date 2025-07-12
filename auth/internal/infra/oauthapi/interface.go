package oauthapi

import oauthmodels "mandacode.com/accounts/auth/internal/models/oauth"

type OAuthAPI interface {
	// GetAccessToken retrieves an access token using the provided authorization code.
	//
	// Parameters:
	//   - code: The authorization code received from the OAuth provider.
	//
	// Returns:
	//   - A string representing the access token.
	//   - An error if the token retrieval fails.
	GetAccessToken(code string) (string, error)

	// GetLoginURL returns the URL to redirect the user for OAuth login.
	GetLoginURL() string

	// GetUserInfo retrieves user information using the access token.
	//
	// Parameters:
	//   - accessToken: The access token obtained from the OAuth provider.
	GetUserInfo(accessToken string) (*oauthmodels.UserInfo, error)
}
