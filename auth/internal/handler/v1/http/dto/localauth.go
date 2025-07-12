package handlerv1dto

type LocalLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type IssueCodeResponse struct {
	Code   string `json:"code"`
	UserID string `json:"user_id"`
}

type LocalSignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

