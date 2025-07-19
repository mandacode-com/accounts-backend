package signupinfradto

type OAuthSignupResponse struct {
	UserID    string `json:"user_id"`
	Provider  string `json:"provider"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}
