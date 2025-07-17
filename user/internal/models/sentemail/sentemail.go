package sentemail

import (
	"time"

	"github.com/google/uuid"
	"mandacode.com/accounts/user/ent"
)

type SentEmail struct {
	ID     uuid.UUID `json:"id"`      // Unique identifier for the sent email record. This is a UUID that is generated when the record is created.
	UserID uuid.UUID `json:"user_id"` // Unique identifier for the user associated with the sent email. This is a UUID that is generated when the user is created.
	Email  string    `json:"email"`   // Email address to which the email was sent.
	SentAt time.Time `json:"sent_at"` // Timestamp when the email was sent.
}

// FromEnt converts an ent SentEmail entity to a SentEmail model.
func FromEnt(ent *ent.SentEmail) *SentEmail {
	return &SentEmail{
		ID:     ent.ID,
		UserID: ent.UserID,
		Email:  ent.Email,
		SentAt: ent.SentAt,
	}
}
