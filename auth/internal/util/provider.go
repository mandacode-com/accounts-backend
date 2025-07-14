package util

import (
	providerv1 "github.com/mandacode-com/accounts-proto/go/provider/v1"
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

func FromProtoToEnt(provider providerv1.ProviderType) (authaccount.Provider, error) {
	switch provider {
	case providerv1.ProviderType_PROVIDER_TYPE_GOOGLE:
		return authaccount.ProviderGoogle, nil
	case providerv1.ProviderType_PROVIDER_TYPE_KAKAO:
		return authaccount.ProviderKakao, nil
	case providerv1.ProviderType_PROVIDER_TYPE_NAVER:
		return authaccount.ProviderNaver, nil
	default:
		return "", errors.New("unsupported provider", "UnsupportedProvider", errcode.ErrInvalidInput)
	}
}
