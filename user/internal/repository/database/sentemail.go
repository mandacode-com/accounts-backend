package dbrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"mandacode.com/accounts/user/ent"
	"mandacode.com/accounts/user/ent/sentemail"
)

type SentEmailRepository struct {
	client *ent.Client
}

// NewSentEmailRepository creates a new SentEmailRepository with the provided database client.
func NewSentEmailRepository(client *ent.Client) *SentEmailRepository {
	return &SentEmailRepository{
		client: client,
	}
}

// GetSentEmailByID retrieves a sent email by its ID.
func (r *SentEmailRepository) GetSentEmailByID(ctx context.Context, id uuid.UUID) (*ent.SentEmail, error) {
	sentEmail, err := r.client.SentEmail.Query().
		Where(sentemail.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("Sent email not found", "NotFound", errcode.ErrNotFound)
		}
		return nil, errors.New(err.Error(), "Failed to get SentEmail by ID", errcode.ErrInternalFailure)
	}
	return sentEmail, nil
}

// GetSentEmailsByUserID retrieves a sent email by the user ID.
func (r *SentEmailRepository) GetSentEmailsByUserID(ctx context.Context, userID uuid.UUID) ([]*ent.SentEmail, error) {
	sentEmails, err := r.client.SentEmail.Query().
		Where(sentemail.UserIDEQ(userID)).
		All(ctx)
	if err != nil {
		return nil, errors.New(err.Error(), "Failed to get SentEmails by User ID", errcode.ErrInternalFailure)
	}
	return sentEmails, nil
}

// GetSentEmailsByEmail retrieves a sent email by the email address.
func (r *SentEmailRepository) GetSentEmailsByEmail(ctx context.Context, email string) ([]*ent.SentEmail, error) {
	sentEmails, err := r.client.SentEmail.Query().
		Where(sentemail.EmailEQ(email)).
		All(ctx)
	if err != nil {
		return nil, errors.New(err.Error(), "Failed to get SentEmails by Email", errcode.ErrInternalFailure)
	}
	return sentEmails, nil
}

// GetSentEmailsByUserDuration retrieves sent emails for a user within a specified duration.
func (r *SentEmailRepository) GetSentEmailsByUserDuration(ctx context.Context, userID uuid.UUID, duration time.Duration) ([]*ent.SentEmail, error) {
	sentEmails, err := r.client.SentEmail.Query().
		Where(sentemail.UserIDEQ(userID)).
		Where(sentemail.SentAtGTE(time.Now().Add(-duration))).
		All(ctx)
	if err != nil {
		return nil, errors.New(err.Error(), "Failed to get User SentEmails by Duration", errcode.ErrInternalFailure)
	}
	return sentEmails, nil
}

// GetSentEmailsByEmailDuration retrieves sent emails for an email address within a specified duration.
func (r *SentEmailRepository) GetSentEmailsByEmailDuration(ctx context.Context, email string, duration time.Duration) ([]*ent.SentEmail, error) {
	sentEmails, err := r.client.SentEmail.Query().
		Where(sentemail.EmailEQ(email)).
		Where(sentemail.SentAtGTE(time.Now().Add(-duration))).
		All(ctx)
	if err != nil {
		return nil, errors.New(err.Error(), "Failed to get Email SentEmails by Duration", errcode.ErrInternalFailure)
	}
	return sentEmails, nil
}

// GetSentEmailNumberByUserDuration retrieves the number of sent emails for a user within a specified duration.
func (r *SentEmailRepository) GetSentEmailNumberByUserDuration(ctx context.Context, userID uuid.UUID, duration time.Duration) (int, error) {
	count, err := r.client.SentEmail.Query().
		Where(sentemail.UserIDEQ(userID)).
		Where(sentemail.SentAtGTE(time.Now().Add(-duration))).
		Count(ctx)
	if err != nil {
		return 0, errors.New(err.Error(), "Failed to get User SentEmail count by Duration", errcode.ErrInternalFailure)
	}
	return count, nil
}

// GetSentEmailNumberByEmailDuration retrieves the number of sent emails for an email address within a specified duration.
func (r *SentEmailRepository) GetSentEmailNumberByEmailDuration(ctx context.Context, email string, duration time.Duration) (int, error) {
	count, err := r.client.SentEmail.Query().
		Where(sentemail.EmailEQ(email)).
		Where(sentemail.SentAtGTE(time.Now().Add(-duration))).
		Count(ctx)
	if err != nil {
		return 0, errors.New(err.Error(), "Failed to get Email SentEmail count by Duration", errcode.ErrInternalFailure)
	}
	return count, nil
}

// CreateSentEmail creates a new sent email record with the provided details.
func (r *SentEmailRepository) CreateSentEmail(ctx context.Context, userID uuid.UUID, email string) (*ent.SentEmail, error) {
	create := r.client.SentEmail.Create().
		SetUserID(userID).
		SetEmail(email).
		SetSentAt(time.Now())

	sentEmail, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, errors.New("Sent email already exists", "Conflict", errcode.ErrConflict)
		}
		return nil, errors.New(err.Error(), "Failed to create SentEmail", errcode.ErrInternalFailure)
	}
	return sentEmail, nil
}

// DeleteSentEmail deletes a sent email by its ID.
func (r *SentEmailRepository) DeleteSentEmail(ctx context.Context, id uuid.UUID) error {
	err := r.client.SentEmail.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("Sent email not found", "NotFound", errcode.ErrNotFound)
		}
		return errors.New(err.Error(), "Failed to delete SentEmail", errcode.ErrInternalFailure)
	}
	return nil
}
