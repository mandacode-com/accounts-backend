package schema

import (
	"context"
	"errors"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	gen "mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/hook"
)

// AuthAccount holds the schema definition for the AuthAccount entity.
type AuthAccount struct {
	ent.Schema
}

// Fields of the AuthAccount.
func (AuthAccount) Fields() []ent.Field {
	return []ent.Field{
		// Internal PK
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Unique().
			Default(uuid.New).
			Comment("The unique identifier for the authentication account"),

		// User ID
		field.UUID("user_id", uuid.UUID{}).
			Comment("The unique identifier for the user associated with this authentication account"),

		// Provider
		field.Enum("provider").
			Values("local", "google", "kakao", "naver", "apple").
			Comment("The OAuth provider used for authentication"),

		// ProviderID
		field.String("provider_id").
			Optional().
			Nillable().
			Comment("The unique identifier provided by the OAuth provider for the user"),

		// IsVerified
		field.Bool("is_verified").
			Default(false).
			Comment("Indicates if the authentication account has verified the email address"),

		// Email
		field.String("email").
			NotEmpty().
			Comment("The email address associated with the authentication account"),

		// PasswordHash
		field.String("password_hash").
			Optional().
			Nillable().
			Comment("The hashed password for the local authentication, if applicable"),

		// CreatedAt
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("The time when the authentication account was created"),

		// UpdatedAt
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("The time when the authentication account was last updated"),
	}
}

// Indexes of the AuthAccount.
func (AuthAccount) Indexes() []ent.Index {
	return []ent.Index{
		// Unique index for user_id and provider
		index.Fields("user_id", "provider").Unique(),
		// Unique index for provider and provider_id
		index.Fields("provider", "provider_id").Unique(),
		// Unique index for email and provider
		index.Fields("email", "provider").Unique(),
	}
}

// Edges of the AuthAccount.
func (AuthAccount) Edges() []ent.Edge {
	return nil
}

// Hooks of the AuthAccount.
func (AuthAccount) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(validatePasswordByProvider(), ent.OpCreate|ent.OpUpdate),
		hook.On(validateProviderID(), ent.OpCreate|ent.OpUpdate),
	}
}

// validatePasswordByProvider enforces provider-specific rules for password_hash.
func validatePasswordByProvider() ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return hook.AuthAccountFunc(func(ctx context.Context, m *gen.AuthAccountMutation) (ent.Value, error) {
			provider, ok := m.Provider()
			if !ok {
				return nil, errors.New("provider must be set")
			}

			pw, hasPw := m.PasswordHash()

			if provider == "local" {
				if !hasPw || pw == "" {
					return nil, errors.New("password_hash is required for local provider")
				}
			} else {
				if hasPw && pw != "" {
					return nil, errors.New("password_hash must not be set for non-local providers")
				}
			}

			return next.Mutate(ctx, m)
		})
	}
}

// validateProviderID checks if provider_id is set for OAuth providers.
func validateProviderID() ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return hook.AuthAccountFunc(func(ctx context.Context, m *gen.AuthAccountMutation) (ent.Value, error) {
			provider, ok := m.Provider()
			if !ok {
				return nil, errors.New("provider must be set")
			}

			providerID, hasProviderID := m.ProviderID()

			if provider != "local" && (!hasProviderID || providerID == "") {
				return nil, errors.New("provider_id is required for non-local providers")
			}

			return next.Mutate(ctx, m)
		})
	}
}
