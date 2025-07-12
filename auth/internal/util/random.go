package util

import (
	"crypto/rand"
	"encoding/hex"
)

type RandomGenerator struct {
	CodeLength int
}

func NewRandomGenerator(codeLength int) *RandomGenerator {
	return &RandomGenerator{
		CodeLength: codeLength,
	}
}

// GenerateSecureRandomCode generates a secure random code of the specified length.
func (rg *RandomGenerator) GenerateSecureRandomCode() (string, error) {
	bytes := make([]byte, rg.CodeLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
