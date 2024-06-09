package subscriber

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/omran95/chat-app/pkg/room"
)

var MessagePubTopic = "chat.msg.pub"

type MessagePublisher interface {
	PublishToSubscribers(ctx context.Context, subscribers []string, message room.Message) error
}

type MessagePublisherImpl struct {
	publisher message.Publisher
}

func NewMessagePublisher(publisher message.Publisher) *MessagePublisherImpl {
	return &MessagePublisherImpl{publisher}
}

func (msgPub *MessagePublisherImpl) PublishToSubscribers(ctx context.Context, subscribers []string, msg room.Message) error {
	fmt.Print(subscribers)
	for _, subscriber := range subscribers {
		err := msgPub.publisher.Publish(subscriber, message.NewMessage(watermill.NewULID(), msg.Encode()))
		if err != nil {
			return err
		}
	}
	return nil
}
