package infra_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/google/uuid"
	tokengendomain "mandacode.com/accounts/token/internal/domain/infra/token"
	tokengen "mandacode.com/accounts/token/internal/infra/token"
)

type MockTokenGenerator struct {
	svc tokengendomain.TokenGenerator
}

func (m *MockTokenGenerator) Setup(t *testing.T) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	m.svc, err = tokengen.NewTokenGenerator(priv, time.Second)
	if err != nil {
		t.Fatalf("failed to create token generator: %v", err)
	}
	if m.svc == nil {
		t.Fatal("token generator is nil")
	}
}

func (m *MockTokenGenerator) Teardown() {
	m.svc = nil
}

func TestTokenGenerator_GenerateToken(t *testing.T) {
	mockGen := &MockTokenGenerator{}
	mockGen.Setup(t)
	defer mockGen.Teardown()

	mockUserID := uuid.New().String()

	claims := map[string]string{"sub": mockUserID}

	t.Run("GenerateToken_Success", func(t *testing.T) {
		token, exp, err := mockGen.svc.GenerateToken(claims)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if token == "" {
			t.Error("expected a non-empty token")
		}
		if exp <= 0 {
			t.Error("expected a positive expiration time")
		}
	})
}
