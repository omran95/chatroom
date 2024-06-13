package subscriber

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/omran95/chatroom/pkg/room"
)

type MessagePublisher interface {
	PublishToSubscribers(ctx context.Context, subscribers map[string]struct{}, message room.Message) error
}

type MessagePublisherImpl struct {
	publisher message.Publisher
}

func NewMessagePublisher(publisher message.Publisher) *MessagePublisherImpl {
	return &MessagePublisherImpl{publisher}
}

func (msgPub *MessagePublisherImpl) PublishToSubscribers(ctx context.Context, subscribers map[string]struct{}, msg room.Message) error {
	for subscriber, _ := range subscribers {
		err := msgPub.publisher.Publish(subscriber, message.NewMessage(watermill.NewULID(), msg.Encode()))
		if err != nil {
			return err
		}
	}
	return nil
}
