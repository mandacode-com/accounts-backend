package usereventrepo

import (
	"context"

	"github.com/google/uuid"
	usereventv1 "github.com/mandacode-com/accounts-proto/go/user/event/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type UserEventEmitter struct {
	writer *kafka.Writer
}

// NewUserEventEmitter creates a new UserEventEmitter with the provided Kafka writer.
func NewUserEventEmitter(writer *kafka.Writer) *UserEventEmitter {
	return &UserEventEmitter{
		writer: writer,
	}
}

// EmitUserDeletedEvent emits a user deleted event to Kafka.
func (e *UserEventEmitter) EmitUserDeletedEvent(ctx context.Context, userID uuid.UUID) error {
	event := &usereventv1.UserEvent{
		EventType: usereventv1.EventType_USER_DELETED,
		UserId:    userID.String(),
		SyncCode:  nil,
	}

	// Marshal the event to protobuf bytes
	data, err := proto.Marshal(event)
	if err != nil {
		return errors.New(err.Error(), "Failed to marshal user deleted event", errcode.ErrInternalFailure)
	}

	// Create a message to send to Kafka
	message := kafka.Message{
		Key:   []byte(userID.String()),
		Value: data,
	}

	if err := e.writer.WriteMessages(ctx, message); err != nil {
		return errors.New(err.Error(), "Failed to write user deleted event to Kafka", errcode.ErrInternalFailure)
	}
	return nil
}

// EmitUserArchivedEvent emits a user archived event to Kafka.
func (e *UserEventEmitter) EmitUserArchivedEvent(ctx context.Context, userID uuid.UUID, syncCode string) error {
	event := &usereventv1.UserEvent{
		EventType: usereventv1.EventType_USER_ARCHIVED,
		UserId:    userID.String(),
		SyncCode:  &syncCode,
	}

	// Marshal the event to protobuf bytes
	data, err := proto.Marshal(event)
	if err != nil {
		return errors.New(err.Error(), "Failed to marshal user archived event", errcode.ErrInternalFailure)
	}

	// Create a message to send to Kafka
	message := kafka.Message{
		Key:   []byte(userID.String()),
		Value: data,
	}

	if err := e.writer.WriteMessages(ctx, message); err != nil {
		return errors.New(err.Error(), "Failed to write user archived event to Kafka", errcode.ErrInternalFailure)
	}
	return nil
}

// EmitUserRestoredEvent emits a user restored event to Kafka.
func (e *UserEventEmitter) EmitUserRestoredEvent(ctx context.Context, userID uuid.UUID, syncCode string) error {
	event := &usereventv1.UserEvent{
		EventType: usereventv1.EventType_USER_RESTORED,
		UserId:    userID.String(),
		SyncCode:  &syncCode,
	}

	// Marshal the event to protobuf bytes
	data, err := proto.Marshal(event)
	if err != nil {
		return errors.New(err.Error(), "Failed to marshal user restored event", errcode.ErrInternalFailure)
	}

	// Create a message to send to Kafka
	message := kafka.Message{
		Key:   []byte(userID.String()),
		Value: data,
	}

	if err := e.writer.WriteMessages(ctx, message); err != nil {
		return errors.New(err.Error(), "Failed to write user restored event to Kafka", errcode.ErrInternalFailure)
	}
	return nil
}

// EmitUserBlockedEvent emits a user blocked event to Kafka.
func (e *UserEventEmitter) EmitUserBlockedEvent(ctx context.Context, userID uuid.UUID, syncCode string) error {
	event := &usereventv1.UserEvent{
		EventType: usereventv1.EventType_USER_BLOCKED,
		UserId:    userID.String(),
		SyncCode:  &syncCode,
	}

	// Marshal the event to protobuf bytes
	data, err := proto.Marshal(event)
	if err != nil {
		return errors.New(err.Error(), "Failed to marshal user blocked event", errcode.ErrInternalFailure)
	}

	// Create a message to send to Kafka
	message := kafka.Message{
		Key:   []byte(userID.String()),
		Value: data,
	}

	if err := e.writer.WriteMessages(ctx, message); err != nil {
		return errors.New(err.Error(), "Failed to write user blocked event to Kafka", errcode.ErrInternalFailure)
	}
	return nil
}

// EmitUserUnblockedEvent emits a user unblocked event to Kafka.
func (e *UserEventEmitter) EmitUserUnblockedEvent(ctx context.Context, userID uuid.UUID, syncCode string) error {
	event := &usereventv1.UserEvent{
		EventType: usereventv1.EventType_USER_UNBLOCKED,
		UserId:    userID.String(),
		SyncCode:  &syncCode,
	}

	// Marshal the event to protobuf bytes
	data, err := proto.Marshal(event)
	if err != nil {
		return errors.New(err.Error(), "Failed to marshal user unblocked event", errcode.ErrInternalFailure)
	}

	// Create a message to send to Kafka
	message := kafka.Message{
		Key:   []byte(userID.String()),
		Value: data,
	}

	if err := e.writer.WriteMessages(ctx, message); err != nil {
		return errors.New(err.Error(), "Failed to write user unblocked event to Kafka", errcode.ErrInternalFailure)
	}
	return nil
}
