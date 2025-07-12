package config

import (
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
)

type OAuthProviderConfig struct {
	ClientID     string `validate:"required"`
	ClientSecret string `validate:"required"`
	RedirectURL  string `validate:"required,url"`
}

type KafkaWriterConfig struct {
	Address string `validate:"required"`
	Topic   string `validate:"required"`
}

type KafkaReaderConfig struct {
	Address string `validate:"required"`
	Topic   string `validate:"required"`
	GroupID string `validate:"required"`
}

type RedisStoreConfig struct {
	Address  string        `validate:"required"`
	Password string        `validate:"omitempty"`
	DB       int           `validate:"min=0,max=15"`
	Prefix   string        `validate:"omitempty"`
	HashKey  string        `validate:"required"`
	Timeout  time.Duration `validate:"omitempty,min=1"`
}

type Config struct {
	Env              string              `validate:"required,oneof=dev prod"`
	Port             int                 `validate:"required,min=1,max=65535"`
	TokenServiceAddr string              `validate:"required"`
	DatabaseURL      string              `validate:"required"`
	VerifyEmailURL   string              `validate:"required,url"`
	LoginCodeStore   RedisStoreConfig    `validate:"required"`
	EmailCodeStore   RedisStoreConfig    `validate:"required"` // Store for email verification codes
	SessionStore     RedisStoreConfig    `validate:"required"`
	MailWriter       KafkaWriterConfig   `validate:"required"`
	UserEventReader  KafkaReaderConfig   `validate:"required"`
	GoogleOAuth      OAuthProviderConfig `validate:"required"`
	NaverOAuth       OAuthProviderConfig `validate:"required"`
	KakaoOAuth       OAuthProviderConfig `validate:"required"`
}

// LoadConfig loads env vars from .env (if exists) and returns structured config
func LoadConfig(validator *validator.Validate) (*Config, error) {
	if os.Getenv("ENV") != "prod" {
		_ = godotenv.Load()
	}

	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, err
	}
	sessionStoreDB, err := strconv.Atoi(getEnv("SESSION_STORE_DB", "0"))
	if err != nil {
		return nil, err
	}
	codeStoreDB, err := strconv.Atoi(getEnv("CODE_STORE_DB", "0"))
	if err != nil {
		return nil, err
	}
	loginCodeTTL, err := time.ParseDuration(getEnv("LOGIN_CODE_TTL", "5m"))
	if err != nil {
		return nil, errors.New("Invalid LOGIN_CODE_TTL format", "Failed to parse login code TTL", errcode.ErrInvalidInput)
	}
	emailCodeTTL, err := time.ParseDuration(getEnv("EMAIL_CODE_TTL", "1h"))
	if err != nil {
		return nil, errors.New("Invalid EMAIL_CODE_TTL format", "Failed to parse email code TTL", errcode.ErrInvalidInput)
	}

	config := &Config{
		Env:              getEnv("ENV", "dev"),
		Port:             port,
		TokenServiceAddr: getEnv("TOKEN_SERVICE_ADDR", ""),
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		VerifyEmailURL:   getEnv("VERIFY_EMAIL_URL", ""),
		LoginCodeStore: RedisStoreConfig{
			Address:  getEnv("LOGIN_CODE_STORE_ADDRESS", ""),
			Password: getEnv("LOGIN_CODE_STORE_PASSWORD", ""),
			DB:       codeStoreDB,
			Prefix:   getEnv("LOGIN_CODE_STORE_PREFIX", "login_code:"),
			HashKey:  getEnv("LOGIN_CODE_STORE_HASH_KEY", "default_login_code_hash_key"),
			Timeout:  loginCodeTTL,
		},
		EmailCodeStore: RedisStoreConfig{
			Address:  getEnv("EMAIL_CODE_STORE_ADDRESS", ""),
			Password: getEnv("EMAIL_CODE_STORE_PASSWORD", ""),
			DB:       codeStoreDB,
			Prefix:   getEnv("EMAIL_CODE_STORE_PREFIX", "email_code:"),
			HashKey:  getEnv("EMAIL_CODE_STORE_HASH_KEY", "default_email_code_hash_key"),
			Timeout:  emailCodeTTL,
		},
		SessionStore: RedisStoreConfig{
			Address:  getEnv("SESSION_STORE_ADDRESS", ""),
			Password: getEnv("SESSION_STORE_PASSWORD", ""),
			DB:       sessionStoreDB,
			Prefix:   getEnv("SESSION_STORE_PREFIX", "session:"),
			HashKey:  getEnv("SESSION_STORE_HASH_KEY", "default_session_hash_key"),
		},
		MailWriter: KafkaWriterConfig{
			Address: getEnv("MAIL_WRITER_ADDRESS", ""),
			Topic:   getEnv("MAIL_WRITER_TOPIC", "mail"),
		},
		UserEventReader: KafkaReaderConfig{
			Address: getEnv("USER_EVENT_READER_ADDRESS", ""),
			Topic:   getEnv("USER_EVENT_READER_TOPIC", "user_event"),
			GroupID: getEnv("USER_EVENT_READER_GROUP_ID", "user_event_group"),
		},
		GoogleOAuth: OAuthProviderConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
		},
		NaverOAuth: OAuthProviderConfig{
			ClientID:     getEnv("NAVER_CLIENT_ID", ""),
			ClientSecret: getEnv("NAVER_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("NAVER_REDIRECT_URL", ""),
		},
		KakaoOAuth: OAuthProviderConfig{
			ClientID:     getEnv("KAKAO_CLIENT_ID", ""),
			ClientSecret: getEnv("KAKAO_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("KAKAO_REDIRECT_URL", ""),
		},
	}

	if err := validator.Struct(config); err != nil {
		return nil, errors.New(err.Error(), "Invalid configuration", errcode.ErrInvalidInput)
	}
	return config, nil
}

// getEnv returns env value or fallback
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
