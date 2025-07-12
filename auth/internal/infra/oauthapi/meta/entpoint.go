package oauthapimeta

const (
	GoogleTokenEndpoint = "https://oauth2.googleapis.com/token"
	NaverTokenEndpoint  = "https://nid.naver.com/oauth2.0/token"
	KakaoTokenEndpoint  = "https://kauth.kakao.com/oauth/token"
)

const (
	GoogleGrantType = "authorization_code"
	NaverGrantType  = "authorization_code"
	KakaoGrantType  = "authorization_code"
)

const (
	GoogleAuthEndpoint = "https://accounts.google.com/o/oauth2/auth"
	NaverAuthEndpoint  = "https://nid.naver.com/oauth2.0/authorize"
	KakaoAuthEndpoint  = "https://kauth.kakao.com/oauth/authorize"
)

const (
	GoogleUserInfoEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
	KakaoUserInfoEndpoint = "https://kapi.kakao.com/v2/user/me"
	NaverUserInfoEndpoint = "https://openapi.naver.com/v1/nid/me"
)
