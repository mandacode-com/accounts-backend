package dbrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"golang.org/x/crypto/bcrypt"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/authaccount"
	dbmodels "mandacode.com/accounts/auth/internal/models/database"
)

type AuthAccountRepository struct {
	client *ent.Client
}

// CreateLocalAuthAccount creates a new local authentication account.
func (a *AuthAccountRepository) CreateLocalAuthAccount(ctx context.Context, account *dbmodels.CreateLocalAuthAccountInput) (*dbmodels.SecureLocalAuthAccount, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New(err.Error(), "Failed to generate password hash", errcode.ErrInternalFailure)
	}

	create := a.client.AuthAccount.Create().
		SetID(uuid.New()).
		SetUserID(account.UserID).
		SetProvider("local").
		SetEmail(account.Email).
		SetIsVerified(account.IsVerified).
		SetPasswordHash(string(passwordHash))

	authAccount, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, errors.New("AuthAccount already exists", "User Already Exists", errcode.ErrConflict)
		}
		return nil, errors.New(err.Error(), "Failed to create Local AuthAccount", errcode.ErrInternalFailure)
	}

	// return authAccount, nil
	return dbmodels.NewSecureLocalAuthAccount(authAccount), nil
}

// CreateLocalAuthAccount creates a new oauth authentication account.
func (a *AuthAccountRepository) CreateOAuthAuthAccount(ctx context.Context, account *dbmodels.CreateOAuthAuthAccountInput) (*dbmodels.SecureOAuthAuthAccount, error) {
	create := a.client.AuthAccount.Create().
		SetID(uuid.New()).
		SetUserID(account.UserID).
		SetProvider(account.Provider).
		SetProviderID(account.ProviderID).
		SetEmail(account.Email).
		SetIsVerified(account.IsVerified)

	authAccount, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, errors.New("AuthAccount already exists", "User Already Exists", errcode.ErrConflict)
		}
		return nil, errors.New(err.Error(), "Failed to create OAuth AuthAccount", errcode.ErrInternalFailure)
	}

	return dbmodels.NewSecureOAuthAuthAccount(authAccount), nil
}

// GetAuthAccountByID retrieves an authentication account by ID.
func (a *AuthAccountRepository) GetAuthAccountByID(ctx context.Context, id uuid.UUID) (*dbmodels.SecureAuthAccount, error) {
	authAccount, err := a.client.AuthAccount.Query().
		Where(authaccount.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found", "AuthAccount Not Found", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to find AuthAccount by ID", errcode.ErrInternalFailure)
	}

	return dbmodels.NewSecureAuthAccount(authAccount), nil
}

// GetAuthAccountsByUserID retrieves an authentication account by user ID.
func (a *AuthAccountRepository) GetAuthAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]*dbmodels.SecureAuthAccount, error) {
	authAccounts, err := a.client.AuthAccount.Query().
		Where(authaccount.UserID(userID)).
		All(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found", "AuthAccount Not Found", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to find AuthAccounts by UserID", errcode.ErrInternalFailure)
	}

	var secureAccounts []*dbmodels.SecureAuthAccount
	for _, account := range authAccounts {
		secureAccounts = append(secureAccounts, dbmodels.NewSecureAuthAccount(account))
	}

	return secureAccounts, nil
}

// GetLocalAuthAccountByUserID retrieves a local authentication account by user ID.
func (a *AuthAccountRepository) GetLocalAuthAccountByUserID(ctx context.Context, userID uuid.UUID) (*dbmodels.SecureLocalAuthAccount, error) {
	authAccount, err := a.client.AuthAccount.Query().
		Where(authaccount.And(
			authaccount.UserIDEQ(userID),
			authaccount.ProviderEQ(authaccount.ProviderLocal),
		)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found", "AuthAccount Not Found", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to find AuthAccount by UserID", errcode.ErrInternalFailure)
	}

	return dbmodels.NewSecureLocalAuthAccount(authAccount), nil
}

// GetLocalAuthAccountByEmail retrieves a local authentication account by email.
func (a *AuthAccountRepository) GetLocalAuthAccountByEmail(ctx context.Context, email string) (*dbmodels.SecureLocalAuthAccount, error) {
	authAccount, err := a.client.AuthAccount.Query().
		Where(authaccount.And(
			authaccount.Email(email),
			authaccount.ProviderEQ(authaccount.ProviderLocal),
		)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found", "AuthAccount Not Found", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to find AuthAccount by Email", errcode.ErrInternalFailure)
	}

	return dbmodels.NewSecureLocalAuthAccount(authAccount), nil
}

// GetOAuthAuthAccountByUserID retrieves an OAuth authentication account by user ID.
func (a *AuthAccountRepository) GetOAuthAuthAccountByUserID(ctx context.Context, userID uuid.UUID, provider authaccount.Provider) (*dbmodels.SecureOAuthAuthAccount, error) {
	authAccount, err := a.client.AuthAccount.Query().
		Where(authaccount.And(
			authaccount.UserID(userID),
			authaccount.ProviderEQ(provider),
		)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found", "AuthAccount Not Found", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to find AuthAccount by UserID and Provider", errcode.ErrInternalFailure)
	}
	return dbmodels.NewSecureOAuthAuthAccount(authAccount), nil
}

// GetOAuthAccountByProviderAndProviderID retrieves an OAuth authentication account by provider and provider ID.
func (a *AuthAccountRepository) GetOAuthAccountByProviderAndProviderID(ctx context.Context, provider authaccount.Provider, providerID string) (*dbmodels.SecureOAuthAuthAccount, error) {
	if provider == authaccount.ProviderLocal {
		return nil, errors.New("Invalid provider", "Provider cannot be 'local' for OAuth accounts", errcode.ErrInvalidInput)
	}

	authAccount, err := a.client.AuthAccount.Query().
		Where(authaccount.And(
			authaccount.ProviderEQ(provider),
			authaccount.ProviderID(providerID),
		)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found", "AuthAccount Not Found", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to find AuthAccount by Provider and ProviderID", errcode.ErrInternalFailure)
	}

	return dbmodels.NewSecureOAuthAuthAccount(authAccount), nil
}

// SetPasswordHash sets the password hash for a local authentication account.
func (a *AuthAccountRepository) SetPasswordHash(ctx context.Context, userID uuid.UUID, password string) (*dbmodels.SecureLocalAuthAccount, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New(err.Error(), "Failed to generate password hash", errcode.ErrInternalFailure)
	}

	localAccount, err := a.client.AuthAccount.Query().
		Where(authaccount.And(
			authaccount.UserID(userID),
			authaccount.ProviderEQ(authaccount.ProviderLocal),
		)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found", "AuthAccount Not Found", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to find Local AuthAccount", errcode.ErrInternalFailure)
	}

	update := localAccount.Update().
		SetPasswordHash(string(passwordHash))
	authAccount, err := update.Save(ctx)
	if err != nil {
		return nil, errors.New(err.Error(), "Failed to update Local AuthAccount password hash", errcode.ErrInternalFailure)
	}

	return dbmodels.NewSecureLocalAuthAccount(authAccount), nil
}

// ComparePassword compares the provided password with the stored password hash.
//
// Parameters:
//   - ctx: The context for the operation.
//   - email: The email of the local authentication account.
//   - password: The password to compare.
//
// Returns:
//   - bool: true if the password matches, false otherwise.
//   - uuid.UUID: The user ID associated with the account.
//   - error: An error if the operation fails, nil otherwise.
func (a *AuthAccountRepository) ComparePassword(ctx context.Context, email string, password string) (bool, uuid.UUID, error) {
	localAccount, err := a.client.AuthAccount.Query().
		Where(authaccount.And(
			authaccount.Email(email),
			authaccount.ProviderEQ(authaccount.ProviderLocal),
		)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return false, uuid.Nil, nil
		}
		return false, uuid.Nil, errors.New(err.Error(), "Internal Error", errcode.ErrInternalFailure)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*localAccount.PasswordHash), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, uuid.Nil, nil // Password does not match, return user ID for further processing
		}
		return false, uuid.Nil, errors.New(err.Error(), "Internal Error", errcode.ErrInternalFailure)
	}

	return true, localAccount.UserID, nil // Password matches, return user ID
}

// DeleteAuthAccountByUserID deletes an authentication account by user ID.
func (a *AuthAccountRepository) DeleteAuthAccountByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := a.client.AuthAccount.Delete().
		Where(authaccount.UserID(userID)).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("AuthAccount not found", "AuthAccount Not Found", errcode.ErrNotFound)
		}
		return errors.New(err.Error(), "Failed to delete AuthAccount by UserID", errcode.ErrInternalFailure)
	}

	return nil
}

// DeleteAuthAccountByID deletes an authentication account by ID.
func (a *AuthAccountRepository) DeleteAuthAccountByID(ctx context.Context, id uuid.UUID) error {
	_, err := a.client.AuthAccount.Delete().
		Where(authaccount.ID(id)).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("AuthAccount not found", "AuthAccount Not Found", errcode.ErrNotFound)
		}
		return errors.New(err.Error(), "Failed to delete AuthAccount by ID", errcode.ErrInternalFailure)
	}

	return nil
}

// DeleteAuthAccountByUserIDAndProvider deletes an authentication account by user ID and provider.
func (a *AuthAccountRepository) DeleteAuthAccountByUserIDAndProvider(ctx context.Context, userID uuid.UUID, provider authaccount.Provider) error {
	_, err := a.client.AuthAccount.Delete().
		Where(authaccount.And(
			authaccount.UserID(userID),
			authaccount.ProviderEQ(provider),
		)).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("AuthAccount not found", "AuthAccount Not Found", errcode.ErrNotFound)
		}
		return errors.New(err.Error(), "Failed to delete AuthAccount by UserID and Provider", errcode.ErrInternalFailure)
	}

	return nil
}

// SetIsVerifiedByID sets the verification status of a local authentication account by user ID.
func (a *AuthAccountRepository) SetIsVerifiedByID(ctx context.Context, id uuid.UUID, isVerified bool) (*dbmodels.SecureAuthAccount, error) {
	authAccount, err := a.client.AuthAccount.UpdateOneID(id).
		SetIsVerified(isVerified).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found", "AuthAccount Not Found", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to update AuthAccount verification status", errcode.ErrInternalFailure)
	}

	return dbmodels.NewSecureAuthAccount(authAccount), nil
}

// SetIsVerifiedByUserIDAndProvider sets the verification status of an OAuth authentication account by user ID and provider.
func (a *AuthAccountRepository) SetIsVerifiedByUserIDAndProvider(ctx context.Context, userID uuid.UUID, provider authaccount.Provider, isVerified bool) (*dbmodels.SecureAuthAccount, error) {
	authAccount, err := a.client.AuthAccount.Query().
		Where(authaccount.And(
			authaccount.UserID(userID),
			authaccount.ProviderEQ(provider),
		)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found", "AuthAccount Not Found", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to find AuthAccount by UserID and Provider", errcode.ErrInternalFailure)
	}

	update := authAccount.Update().
		SetIsVerified(isVerified)

	authAccountUpdated, err := update.Save(ctx)
	if err != nil {
		return nil, errors.New(err.Error(), "Failed to update AuthAccount verification status", errcode.ErrInternalFailure)
	}

	return dbmodels.NewSecureAuthAccount(authAccountUpdated), nil
}

// UpdateEmailByID updates the email of an authentication account by user ID.
func (a *AuthAccountRepository) UpdateEmailByID(ctx context.Context, id uuid.UUID, newEmail string) (*dbmodels.SecureAuthAccount, error) {
	authAccount, err := a.client.AuthAccount.UpdateOneID(id).
		SetEmail(newEmail).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("AuthAccount not found", "AuthAccount Not Found", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to update AuthAccount email", errcode.ErrInternalFailure)
	}

	return dbmodels.NewSecureAuthAccount(authAccount), nil
}

// UpdateLocalEmailVerificationStatus updates the email verification status of a local authentication account by user ID
func (a *AuthAccountRepository) UpdateLocalEmailVerificationStatus(ctx context.Context, userID uuid.UUID, isVerified bool) error {
	authAccount, err := a.client.AuthAccount.Query().
		Where(authaccount.And(
			authaccount.UserID(userID),
			authaccount.ProviderEQ(authaccount.ProviderLocal),
		)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("AuthAccount not found", "AuthAccount Not Found", errcode.ErrNotFound)
		}
		return errors.New(err.Error(), "Failed to find Local AuthAccount by UserID", errcode.ErrInternalFailure)
	}

	update := authAccount.Update().
		SetIsVerified(isVerified)

	if _, err := update.Save(ctx); err != nil {
		return errors.New(err.Error(), "Failed to update Local AuthAccount email verification status", errcode.ErrInternalFailure)
	}

	return nil
}

// NewAuthAccountRepository creates a new instance of authAccountRepository.
func NewAuthAccountRepository(client *ent.Client) *AuthAccountRepository {
	return &AuthAccountRepository{
		client: client,
	}
}
