package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// SentEmail holds the schema definition for the SentEmail entity.
type SentEmail struct {
	ent.Schema
}

// Fields of the MailSent.
func (SentEmail) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable().
			Unique().
			Comment("Unique identifier for the mail sent record. This is a UUID that is generated when the record is created."),
		field.UUID("user_id", uuid.UUID{}).
			Comment("Unique identifier for the user associated with the mail sent record. This is a UUID that is generated when the user is created."),
		field.String("email").
			NotEmpty().
			Comment("Email address to which the mail was sent. This is a required field and must not be empty."),
		field.Time("sent_at").
			Default(time.Now).
			Immutable().
			Comment("Timestamp when the mail was sent. This is set to the current time when the record is created."),
	}
}

// Edges of the MailSent.
func (SentEmail) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("sent_emails").
			Unique().
			Field("user_id").
			Required().
			Comment("Edge to the User entity. This establishes a relationship between the SentEmail and User entities, linking the sent email to the user who sent it."),
	}
}
