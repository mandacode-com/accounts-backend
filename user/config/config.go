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

type KafkaWriterConfig struct {
	Address string `validate:"required"`
	Topic   string `validate:"required"`
}

type KafkaReaderConfig struct {
	Address string `validate:"required"`
	Topic   string `validate:"required"`
	GroupID string `validate:"required"`
}

type HTTPServerConfig struct {
	Port int `validate:"required,min=1,max=65535"`
}

type GRPCServerConfig struct {
	Port int `validate:"required,min=1,max=65535"`
}

type GRPCClientConfig struct {
	Address string        `validate:"required"`
	Timeout time.Duration `validate:"required,min=1"`
}

type Config struct {
	Env             string            `validate:"required,oneof=dev prod"`
	DatabaseURL     string            `validate:"required"`
	HTTPServer      HTTPServerConfig  `validate:"required"`
	GRPCServer      GRPCServerConfig  `validate:"required"`
	MailWriter      KafkaWriterConfig `validate:"required"`
	UserEventWriter KafkaWriterConfig `validate:"required"`
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
	grpcPort, err := strconv.Atoi(getEnv("GRPC_PORT", "9090"))
	if err != nil {
		return nil, err
	}

	config := &Config{
		Env:              getEnv("ENV", "dev"),
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		MailWriter: KafkaWriterConfig{
			Address: getEnv("MAIL_WRITER_ADDRESS", ""),
			Topic:   getEnv("MAIL_WRITER_TOPIC", ""),
		},
		UserEventWriter: KafkaWriterConfig{
			Address: getEnv("USER_EVENT_WRITER_ADDRESS", ""),
			Topic:   getEnv("USER_EVENT_WRITER_TOPIC", ""),
		},
		HTTPServer: HTTPServerConfig{
			Port: httpPort,
		},
		GRPCServer: GRPCServerConfig{
			Port: grpcPort,
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
