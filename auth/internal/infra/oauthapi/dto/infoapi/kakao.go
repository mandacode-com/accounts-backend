package infoapidto

// RawKakaoUserInfo represents the raw user info structure returned by Kakao OAuth.
type RawKakaoUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"kakao_account.email"`
	Name          string `json:"properties.nickname"`
	EmailValid    bool   `json:"kakao_account.is_email_valid"`
	EmailVerified bool   `json:"kakao_account.is_email_verified"`
}
