package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Profile holds the schema definition for the Profile entity.
type Profile struct {
	ent.Schema
}

// Fields of the Profile.
func (Profile) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("user_id", uuid.UUID{}).
			Immutable().
			Unique().
			Comment("The ID of the user this profile belongs to. This is immutable and should not change."),
		field.String("avatar").
			Optional().
			Comment("The URL of the user's avatar image. This is optional and can be updated."),
		field.String("bio").
			Optional().
			Comment("A short biography of the user. This is optional and can be updated."),
		field.String("location").
			Optional().
			Comment("The user's location. This is optional and can be updated."),
		field.String("nickname").
			Unique().
			Comment("The user's nickname. This is optional and can be updated."),
		field.String("email").
			Optional().
			Comment("The user's email address. This is optional and can be updated."),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("The timestamp when the profile was created. This is immutable and should not change."),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("The timestamp when the profile was last updated. This is automatically set to the current time on update."),
		field.Bool("is_archived").
			Optional().
			Comment("The timestamp when the profile was archived. This is optional and can be set to indicate when the profile was archived."),
		field.Time("archived_at").
			Nillable().
			Optional().
			Comment("The timestamp when the profile was archived. This is optional and can be set to indicate when the profile was archived."),
	}
}

// Indexes of the Profile.
func (Profile) Indexes() []ent.Index {
	return []ent.Index{
		// Index for user_id to ensure uniqueness and improve query performance.
		index.Fields("user_id").
			Unique(),
	}
}

// Edges of the Profile.
func (Profile) Edges() []ent.Edge {
	return nil
}
