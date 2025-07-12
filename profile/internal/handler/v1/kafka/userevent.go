package kafkahandlerv1

import (
	"context"

	"github.com/google/uuid"
	usereventv1 "github.com/mandacode-com/accounts-proto/go/user/event/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	kafka "github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	kafkaserver "mandacode.com/accounts/profile/cmd/server/kafka"
	"mandacode.com/accounts/profile/internal/usecase/system"
)

type UserEventHandler struct {
	// userEvent *userevent.UserEventUsecase
	profile *system.ProfileUsecase
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
		if err := u.profile.DeleteProfile(ctx, userUUID); err != nil {
			return errors.Upgrade(err, "Failed to handle user deleted event", errcode.ErrInternalFailure)
		}
	case usereventv1.EventType_USER_ARCHIVED:
		if err := u.profile.ArchiveProfile(ctx, userUUID); err != nil {
			return errors.Upgrade(err, "Failed to handle user archived event", errcode.ErrInternalFailure)
		}
	case usereventv1.EventType_USER_RESTORED:
		if err := u.profile.RestoreProfile(ctx, userUUID); err != nil {
			return errors.Upgrade(err, "Failed to handle user restored event", errcode.ErrInternalFailure)
		}
	case usereventv1.EventType_USER_BLOCKED:
		return nil
	case usereventv1.EventType_USER_UNBLOCKED:
		return nil
	default:
		return errors.New("unsupported user event type", "User Event Handler Error", errcode.ErrInvalidInput)
	}

	return nil
}

func NewUserEventHandler(profile *system.ProfileUsecase) kafkaserver.KafkaHandler {
	return &UserEventHandler{
		profile: profile,
	}
}
