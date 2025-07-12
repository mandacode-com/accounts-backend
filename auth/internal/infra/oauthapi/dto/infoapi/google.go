package infoapidto

// RawGoogleUserInfo represents the raw user info structure returned by Google OAuth.
type RawGoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	EmailVerified bool   `json:"email_verified"`
}
