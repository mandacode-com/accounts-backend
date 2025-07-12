package coderepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"github.com/redis/go-redis/v9"
	"mandacode.com/accounts/auth/internal/util"
)

type CodeManager struct {
	codeGen   *util.RandomGenerator
	codeTTL   time.Duration
	codeStore *redis.Client
	prefix    string
}

// IssueCode issues a new login code for the given user ID.
//
// Parameters:
//   - ctx: The context for the operation.
//   - userID: The unique identifier of the user.
//
// Returns:
//   - A string representing the issued login code.
//   - An error if the code could not be issued.
func (l *CodeManager) IssueCode(ctx context.Context, userID uuid.UUID) (string, error) {
	code, err := l.codeGen.GenerateSecureRandomCode()
	if err != nil {
		return "", err
	}

	key := l.prefix + code

	err = l.codeStore.Set(ctx, key, userID, l.codeTTL).Err()
	if err != nil {
		return "", err
	}

	return code, nil
}

// ValidateCode validates the provided login code for the given user ID.
//
// Parameters:
//   - ctx: The context for the operation.
//   - userID: The unique identifier of the user.
//   - code: The login code to validate.
//
// Returns:
//   - A boolean indicating whether the code is valid.
//   - An error if the validation fails.
func (l *CodeManager) ValidateCode(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	key := l.prefix + code
	storedUserID, err := l.codeStore.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // Code does not exist
		}
		return false, errors.New(err.Error(), "Failed to get login code from store", errcode.ErrInternalFailure)
	}

	if storedUserID != userID.String() {
		return false, nil // Code exists but does not match user ID
	}

	// Delete the code after successful validation
	err = l.codeStore.Del(ctx, code).Err()
	if err != nil {
		return false, errors.New(err.Error(), "Failed to delete login code from store", errcode.ErrInternalFailure)
	}

	return true, nil // Code is valid and deleted
}

func NewCodeManager(codeGen *util.RandomGenerator, codeTTL time.Duration, codeStore *redis.Client, prefix string) *CodeManager {
	return &CodeManager{
		codeGen:   codeGen,
		codeTTL:   codeTTL,
		codeStore: codeStore,
		prefix:    prefix,
	}
}
