package util

import "crypto/rand"

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type RandomStringGenerator struct {
	length int
}

func NewRandomStringGenerator(length int) *RandomStringGenerator {
	return &RandomStringGenerator{length: length}
}

func (r *RandomStringGenerator) Generate() (string, error) {
	b := make([]byte, r.length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}

	return string(b), nil
}
