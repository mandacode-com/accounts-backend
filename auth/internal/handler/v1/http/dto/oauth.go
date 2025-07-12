package handlerv1dto

type MobileOAuthLoginRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

type OAuthCallbackResponse struct {
	Code   string `json:"code"`
	UserID string `json:"user_id"`
}

