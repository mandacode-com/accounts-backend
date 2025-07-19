package config

import (
	"os"
	"strconv"
	"strings"
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
	Address []string `validate:"required"`
	Topic   string   `validate:"required"`
}

type KafkaReaderConfig struct {
	Brokers []string `validate:"required"`
	Topic   string   `validate:"required"`
	GroupID string   `validate:"required"`
}

type RedisStoreConfig struct {
	Address  string        `validate:"required"`
	Password string        `validate:"omitempty"`
	DB       int           `validate:"min=0,max=15"`
	Prefix   string        `validate:"omitempty"`
	HashKey  string        `validate:"required"`
	Timeout  time.Duration `validate:"omitempty,min=1"`
}

type SessionStoreConfig struct {
	Address     string `validate:"required"`
	Password    string `validate:"omitempty"`
	DB          int    `validate:"min=0,max=15"`
	SessionName string `validate:"required"`
	HashKey     string `validate:"required"`
}

type HTTPServerConfig struct {
	Port int `validate:"required,min=1,max=65535"`
}
type GRPCServerConfig struct {
	Port int `validate:"required,min=1,max=65535"`
}

type GRPCClientConfig struct {
	Address string `validate:"required"`
}
type SignupAPIConfig struct {
	Endpoint string        `validate:"required,url"`
	Timeout  time.Duration `validate:"required,min=1"`
}

type Config struct {
	Env             string              `validate:"required,oneof=dev prod"`
	HTTPServer      HTTPServerConfig    `validate:"required"`
	GRPCServer      GRPCServerConfig    `validate:"required"`
	TokenClient     GRPCClientConfig    `validate:"required"`
	DatabaseURL     string              `validate:"required"`
	LoginCodeStore  RedisStoreConfig    `validate:"required"`
	SessionStore    SessionStoreConfig  `validate:"required"`
	UserEventReader KafkaReaderConfig   `validate:"required"`
	SignupAPI       SignupAPIConfig     `validate:"required"`
	GoogleOAuth     OAuthProviderConfig `validate:"required"`
	NaverOAuth      OAuthProviderConfig `validate:"required"`
	KakaoOAuth      OAuthProviderConfig `validate:"required"`
}

// LoadConfig loads env vars from .env (if exists) and returns structured config
func LoadConfig(validator *validator.Validate) (*Config, error) {
	if os.Getenv("ENV") != "prod" {
		_ = godotenv.Load()
	}

	httpPort, err := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	if err != nil {
		return nil, err
	}
	grpcPort, err := strconv.Atoi(getEnv("GRPC_PORT", "50051"))
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
	signupTimeout, err := time.ParseDuration(getEnv("SIGNUP_API_TIMEOUT", "30s"))
	if err != nil {
		return nil, errors.New("Invalid SIGNUP_API_TIMEOUT format", "Failed to parse signup API timeout", errcode.ErrInvalidInput)
	}

	config := &Config{
		Env: getEnv("ENV", "dev"),
		HTTPServer: HTTPServerConfig{
			Port: httpPort,
		},
		GRPCServer: GRPCServerConfig{
			Port: grpcPort,
		},
		TokenClient: GRPCClientConfig{
			Address: getEnv("TOKEN_CLIENT_ADDR", ""),
		},
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		LoginCodeStore: RedisStoreConfig{
			Address:  getEnv("LOGIN_CODE_STORE_ADDRESS", ""),
			Password: getEnv("LOGIN_CODE_STORE_PASSWORD", ""),
			DB:       codeStoreDB,
			Prefix:   getEnv("LOGIN_CODE_STORE_PREFIX", "login_code:"),
			HashKey:  getEnv("LOGIN_CODE_STORE_HASH_KEY", "default_login_code_hash_key"),
			Timeout:  loginCodeTTL,
		},
		SessionStore: SessionStoreConfig{
			Address:     getEnv("SESSION_STORE_ADDRESS", ""),
			Password:    getEnv("SESSION_STORE_PASSWORD", ""),
			DB:          sessionStoreDB,
			SessionName: getEnv("SESSION_STORE_NAME", "session"),
			HashKey:     getEnv("SESSION_STORE_HASH_KEY", "default_session_hash_key"),
		},
		UserEventReader: KafkaReaderConfig{
			Brokers: strings.Split(getEnv("USER_EVENT_READER_BROKERS", ""), ","),
			Topic:   getEnv("USER_EVENT_READER_TOPIC", "user_event"),
			GroupID: getEnv("USER_EVENT_READER_GROUP_ID", "user_event_group"),
		},
		SignupAPI: SignupAPIConfig{
			Endpoint: getEnv("SIGNUP_API_ENDPOINT", ""),
			Timeout:  signupTimeout,
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
