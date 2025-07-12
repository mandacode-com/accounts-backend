package util

import (
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/auth/ent/authaccount"
)

func ConvertToEnt(provider string) (authaccount.Provider, error) {
	switch provider {
	case "google":
		return authaccount.ProviderGoogle, nil
	case "kakao":
		return authaccount.ProviderKakao, nil
	case "naver":
		return authaccount.ProviderNaver, nil
	default:
		return "", errors.New("unsupported provider", "UnsupportedProvider", errcode.ErrInvalidInput)
	}
}
