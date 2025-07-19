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
	Address string `validate:"required"`
	Topic   string `validate:"required"`
}

type KafkaReaderConfig struct {
	// Address string `validate:"required"`
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

type HTTPServerConfig struct {
	Port      int    `validate:"required,min=1,max=65535"`
	UIDHeader string `validate:"required"`
}

type GRPCServerConfig struct {
	Port int `validate:"required,min=1,max=65535"`
}

type Config struct {
	Env             string            `validate:"required,oneof=dev prod"`
	DatabaseURL     string            `validate:"required"`
	HTTPServer      HTTPServerConfig  `validate:"required"`
	GRPCServer      GRPCServerConfig  `validate:"required"`
	UserEventReader KafkaReaderConfig `validate:"required"`
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
	if err != nil {
		return nil, err
	}

	config := &Config{
		Env: getEnv("ENV", "dev"),
		HTTPServer: HTTPServerConfig{
			Port:      httpPort,
			UIDHeader: getEnv("UID_HEADER_KEY", "X-User-ID"),
		},
		GRPCServer: GRPCServerConfig{
			Port: grpcPort,
		},
		DatabaseURL: getEnv("DATABASE_URL", ""),
		UserEventReader: KafkaReaderConfig{
			Brokers: strings.Split(getEnv("USER_EVENT_READER_BROKERS", ""), ","),
			Topic:   getEnv("USER_EVENT_READER_TOPIC", "user_event"),
			GroupID: getEnv("USER_EVENT_READER_GROUP_ID", "user_event_group"),
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
