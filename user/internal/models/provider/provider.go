package provider

import providerv1 "github.com/mandacode-com/accounts-proto/go/provider/v1"

type ProviderType string

const (
	ProviderGoogle   ProviderType = "google"
	ProviderKakao    ProviderType = "kakao"
	ProviderNaver    ProviderType = "naver"
	ProviderApple    ProviderType = "apple"
	ProviderFacebook ProviderType = "facebook"
	ProviderGithub   ProviderType = "github"
	ProviderTwitter  ProviderType = "twitter"
	ProviderUnknown  ProviderType = "unknown"
)

func ToLocalProvider(provider providerv1.ProviderType) ProviderType {
	switch provider {
	case providerv1.ProviderType_PROVIDER_TYPE_GOOGLE:
		return ProviderGoogle
	case providerv1.ProviderType_PROVIDER_TYPE_KAKAO:
		return ProviderKakao
	case providerv1.ProviderType_PROVIDER_TYPE_NAVER:
		return ProviderNaver
	case providerv1.ProviderType_PROVIDER_TYPE_APPLE:
		return ProviderApple
	case providerv1.ProviderType_PROVIDER_TYPE_FACEBOOK:
		return ProviderFacebook
	case providerv1.ProviderType_PROVIDER_TYPE_GITHUB:
		return ProviderGithub
	default:
		return ProviderUnknown
	}
}

func FromString(provider string) (ProviderType, error) {
	switch provider {
	case "google":
		return ProviderGoogle, nil
	case "kakao":
		return ProviderKakao, nil
	case "naver":
		return ProviderNaver, nil
	case "apple":
		return ProviderApple, nil
	case "facebook":
		return ProviderFacebook, nil
	case "github":
		return ProviderGithub, nil
	default:
		return ProviderUnknown, nil
	}
}

func (p ProviderType) ToProto() providerv1.ProviderType {
	switch p {
	case ProviderGoogle:
		return providerv1.ProviderType_PROVIDER_TYPE_GOOGLE
	case ProviderKakao:
		return providerv1.ProviderType_PROVIDER_TYPE_KAKAO
	case ProviderNaver:
		return providerv1.ProviderType_PROVIDER_TYPE_NAVER
	case ProviderApple:
		return providerv1.ProviderType_PROVIDER_TYPE_APPLE
	case ProviderFacebook:
		return providerv1.ProviderType_PROVIDER_TYPE_FACEBOOK
	case ProviderGithub:
		return providerv1.ProviderType_PROVIDER_TYPE_GITHUB
	default:
		return providerv1.ProviderType_PROVIDER_TYPE_UNSPECIFIED
	}
}
