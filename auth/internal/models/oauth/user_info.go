package oauthmodels

type UserInfo struct {
	ProviderID    string `validate:"required"`
	Email         string `validate:"required,email"`
	Name          string `validate:"required"`
	EmailVerified bool
}

func NewUserInfo(providerID string, email string, name string, emailVerified bool) *UserInfo {
	return &UserInfo{
		ProviderID:    providerID,
		Email:         email,
		Name:          name,
		EmailVerified: emailVerified,
	}
}
