package maileventrepo

import (
	"context"

	mailerv1 "github.com/mandacode-com/accounts-proto/go/mailer/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MailEventEmitter struct {
	writer *kafka.Writer
}

// SendEmailVerificationMail sends an email verification mail to the user.
//
// Parameters:
//   - email: The email address of the user to send the verification mail to.
//   - verificationLink: The link to be included in the email for verification.
func (m *MailEventEmitter) SendEmailVerificationMail(email string, verificationLink string) error {
	event := &mailerv1.EmailVerificationEvent{
		Email:            email,
		VerificationLink: verificationLink,
		EventTime:        timestamppb.Now(),
	}
	// Marshal the event to protobuf bytes
	data, err := proto.Marshal(event)
	if err != nil {
		return errors.New(err.Error(), "Failed to marshal email verification event", errcode.ErrInternalFailure)
	}

	// Create a message to send to Kafka
	message := kafka.Message{
		Key:   []byte(email),
		Value: data,
	}

	return m.writer.WriteMessages(context.Background(), message)
}

// NewMailEventEmitter creates a new Mailer instance with the provided Kafka writer.
func NewMailEventEmitter(writer *kafka.Writer) *MailEventEmitter {
	return &MailEventEmitter{
		writer: writer,
	}
}
