package kafkahandlerv1

import (
	"context"

	"github.com/google/uuid"
	usereventv1 "github.com/mandacode-com/accounts-proto/go/user/event/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	kafka "github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	kafkaserver "mandacode.com/accounts/auth/cmd/server/kafka"
	"mandacode.com/accounts/auth/internal/usecase/userevent"
)

type UserEventHandler struct {
	userEvent *userevent.UserEventUsecase
}

// HandleMessage implements kafkaserver.KafkaHandler.
func (u *UserEventHandler) HandleMessage(ctx context.Context, m kafka.Message) error {
	event := &usereventv1.UserEvent{}
	if err := proto.Unmarshal(m.Value, event); err != nil {
		return errors.Upgrade(err, "Invalid User Event Message", errcode.ErrInvalidInput)
	}

	userUUID, err := uuid.Parse(event.UserId)
	if err != nil {
		return errors.Upgrade(err, "Invalid User ID in User Event", errcode.ErrInvalidInput)
	}

	switch event.EventType {
	case usereventv1.EventType_USER_DELETED:
		if err := u.userEvent.HandleUserDeleted(ctx, userUUID); err != nil {
			return errors.Upgrade(err, "Failed to handle user deleted event", errcode.ErrInternalFailure)
		}
	case usereventv1.EventType_USER_ARCHIVED:
		return nil
	case usereventv1.EventType_USER_RESTORED:
		return nil
	case usereventv1.EventType_USER_BLOCKED:
		return nil
	case usereventv1.EventType_USER_UNBLOCKED:
		return nil
	default:
		return errors.New("unsupported user event type", "User Event Handler Error", errcode.ErrInvalidInput)
	}

	return nil
}

func NewUserEventHandler(userEvent *userevent.UserEventUsecase) kafkaserver.KafkaHandler {
	return &UserEventHandler{
		userEvent: userEvent,
	}
}
