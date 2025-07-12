package infoapidto

// RawNaverUserInfo represents the raw user info structure returned by Naver OAuth.
type RawNaverUserInfo struct {
	ID         string `json:"response.id"`
	Email      string `json:"response.email"`
	Name       string `json:"response.nickname"`
	ResultCode string `json:"resultcode"`
}
