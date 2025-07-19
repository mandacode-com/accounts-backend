package logindto 

import "mandacode.com/accounts/auth/ent/authaccount"

type OAuthLoginInput struct {
	Provider    authaccount.Provider `json:"provider"`
	AccessToken string               `json:"access_token,omitempty"` // Optional, used for OAuth providers that require an access token
	Code        string               `json:"code,omitempty"`         // Optional, used for OAuth providers that require a code exchange
	// Info        models.RequestInfo `json:"info"`
}
