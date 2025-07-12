package util

import (
	"crypto/rsa"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func LoadRSAPublicKeyFromPEM(keyStr string) (*rsa.PublicKey, error) {
	if keyStr == "" {
		return nil, errors.New("public key string cannot be empty")
	}
	// Load the public key from PEM format
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(keyStr))
	if err != nil {
		return nil, errors.New("failed to parse public key: " + err.Error())
	}
	return publicKey, nil
}

func LoadRSAPrivateKeyFromPEM(keyStr string) (*rsa.PrivateKey, error) {
	if keyStr == "" {
		return nil, errors.New("private key string cannot be empty")
	}
	// Load the private key from PEM format
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(keyStr))
	if err != nil {
		return nil, errors.New("failed to parse private key: " + err.Error())
	}
	return privateKey, nil
}
