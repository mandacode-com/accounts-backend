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

type KafkaWriterConfig struct {
	Address []string `validate:"required"`
	Topic   string   `validate:"required"`
}

type HTTPServerConfig struct {
	Port int `validate:"required,min=1,max=65535"`
}

type RedisStoreConfig struct {
	Address  string        `validate:"required"`
	Password string        `validate:"omitempty"`
	DB       int           `validate:"min=0,max=15"`
	Prefix   string        `validate:"omitempty"`
	HashKey  string        `validate:"required"`
	Timeout  time.Duration `validate:"omitempty,min=1"`
}

type GRPCClientConfig struct {
	Address string `validate:"required"`
}

type Config struct {
	Env                   string            `validate:"required,oneof=dev prod"`
	DatabaseURL           string            `validate:"required"`
	HTTPServer            HTTPServerConfig  `validate:"required"`
	UserEventWriter       KafkaWriterConfig `validate:"required"`
	EmailEventWriter       KafkaWriterConfig `validate:"required"`
	EmailCodeStore        RedisStoreConfig  `validate:"required"`
	AuthClient            GRPCClientConfig  `validate:"required"`
	ProfileClient         GRPCClientConfig  `validate:"required"`
	TokenClient           GRPCClientConfig  `validate:"required"`
	EmailVerificationLink string            `validate:"required,url"`
	UserIDHeaderKey       string            `validate:"required"`
	MaxSentEmails         int               `validate:"required,min=1"`
	MaxSentEmailsDuration time.Duration     `validate:"required,min=1"`
}

// LoadConfig loads env vars from .env (if exists) and returns structured config
func LoadConfig(validator *validator.Validate) (*Config, error) {
	if os.Getenv("ENV") != "prod" {
		_ = godotenv.Load()
	}

	httpPort, err := strconv.Atoi(getEnv("HTTP_PORT", ""))
	if err != nil {
		return nil, err
	}

	emailCodeStoreDB, err := strconv.Atoi(getEnv("EMAIL_CODE_STORE_DB", "0"))
	if err != nil {
		return nil, errors.New("Invalid EMAIL_CODE_STORE_DB format", "Failed to parse email code store DB", errcode.ErrInvalidInput)
	}
	emailCodeTTL, err := time.ParseDuration(getEnv("EMAIL_CODE_TTL", "1h"))
	if err != nil {
		return nil, errors.New("Invalid EMAIL_CODE_TTL format", "Failed to parse email code TTL", errcode.ErrInvalidInput)
	}
	maxSentEmails, err := strconv.Atoi(getEnv("MAX_SENT_EMAILS", "5"))
	if err != nil {
		return nil, errors.New("Invalid MAX_SENT_EMAILS format", "Failed to parse max sent emails", errcode.ErrInvalidInput)
	}
	maxSentEmailsDuration, err := time.ParseDuration(getEnv("MAX_SENT_EMAILS_DURATION", "24h"))
	if err != nil {
		return nil, errors.New("Invalid MAX_SENT_EMAILS_DURATION format", "Failed to parse max sent emails duration", errcode.ErrInvalidInput)
	}

	config := &Config{
		Env:         getEnv("ENV", "dev"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		UserEventWriter: KafkaWriterConfig{
			Address: strings.Split(getEnv("USER_EVENT_WRITER_ADDRESS", ""), ","),
			Topic:   getEnv("USER_EVENT_WRITER_TOPIC", ""),
		},
		EmailEventWriter: KafkaWriterConfig{
			Address: strings.Split(getEnv("EMAIL_EVENT_WRITER_ADDRESS", ""), ","),
			Topic:   getEnv("EMAIL_EVENT_WRITER_TOPIC", ""),
		},
		HTTPServer: HTTPServerConfig{
			Port: httpPort,
		},
		EmailCodeStore: RedisStoreConfig{
			Address:  getEnv("EMAIL_CODE_STORE_ADDRESS", ""),
			Password: getEnv("EMAIL_CODE_STORE_PASSWORD", ""),
			DB:       emailCodeStoreDB,
			Prefix:   getEnv("EMAIL_CODE_STORE_PREFIX", "email_code:"),
			HashKey:  getEnv("EMAIL_CODE_STORE_HASH_KEY", "email_code"),
			Timeout:  emailCodeTTL,
		},
		AuthClient: GRPCClientConfig{
			Address: getEnv("AUTH_CLIENT_ADDRESS", ""),
		},
		ProfileClient: GRPCClientConfig{
			Address: getEnv("PROFILE_CLIENT_ADDRESS", ""),
		},
		TokenClient: GRPCClientConfig{
			Address: getEnv("TOKEN_CLIENT_ADDRESS", ""),
		},
		EmailVerificationLink: getEnv("EMAIL_VERIFICATION_LINK", ""),
		UserIDHeaderKey:       getEnv("USER_ID_HEADER_KEY", "X-User-ID"),
		MaxSentEmails:         maxSentEmails,
		MaxSentEmailsDuration: maxSentEmailsDuration,
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
