package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable().
			Unique().
			Comment("Unique identifier for the user. This is a UUID that is generated when the user is created."),
		field.Bool("is_active").
			Default(true).
			Comment("Indicates if the user is active or not. Inactive users cannot log in."),
		field.Bool("is_blocked").
			Default(false).
			Comment("Indicates if the user is blocked. Blocked users cannot log in, but their data is retained for auditing purposes."),
		field.String("sync_code").
			Optional().
			Comment("A code used for synchronizing user data across different systems. This is optional and can be used for integration purposes."),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Timestamp when the user was created. This is set to the current time when the user is created."),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Timestamp when the user was last updated. This is set to the current time whenever the user is updated."),
		field.Bool("is_archived").
			Default(false).
			Comment("Indicates if the user is archived. Archived users are not active but are retained for historical purposes."),
		field.Time("archived_at").
			Optional().
			Nillable().
			Comment("Timestamp when the user was archived. This is set when the user is archived and can be used for auditing purposes."),
		field.Time("delete_after").
			Optional().
			Nillable().
			Comment("Timestamp after which the user will be deleted. This is set when the user is archived and can be used to schedule deletion of the user data."),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
