package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type MailConfig struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required,min=1"`
	Username string `validate:"required"`
	Password string `validate:"required"`
	Sender   string `validate:"required"`
}

type KafkaConfig struct {
	Address string `validate:"required"`
	Topic   string `validate:"required"`
	GroupID string `validate:"required"`
}

type Config struct {
	Env   string      `validate:"required,oneof=dev prod"`
	Mail  MailConfig  `validate:"required"`
	Kafka KafkaConfig `validate:"required"`
}

// LoadConfig loads env vars from .env (if exists) and returns structured config
func LoadConfig(v *validator.Validate) (*Config, error) {
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	mailPort, err := strconv.Atoi(getEnv("MAIL_PORT", "587"))
	if err != nil {
		return nil, err
	}

	mailConfig := MailConfig{
		Host:     getEnv("MAIL_HOST", ""),
		Port:     mailPort,
		Username: getEnv("MAIL_USER", ""),
		Password: getEnv("MAIL_PASS", ""),
		Sender:   getEnv("MAIL_SENDER", ""),
	}

	kafkaConfig := KafkaConfig{
		Address: getEnv("KAFKA_ADDRESS", ""),
		Topic:   getEnv("KAFKA_TOPIC", ""),
		GroupID: getEnv("KAFKA_GROUP_ID", ""),
	}

	config := &Config{
		Env:  getEnv("ENV", "dev"),
		Mail: mailConfig,
		Kafka: kafkaConfig,
	}

	if err := v.Struct(config); err != nil {
		return nil, errors.New("invalid configuration: " + err.Error())
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
