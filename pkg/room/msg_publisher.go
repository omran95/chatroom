package room

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

var MessagePubTopic = "rc.msg.pub"

type MessagePublisher interface {
	PublishMessage(ctx context.Context, message *Message) error
}

type MessagePublisherImpl struct {
	publisher message.Publisher
}

func NewMessagePublisher(publisher message.Publisher) *MessagePublisherImpl {
	return &MessagePublisherImpl{publisher}
}

func (msgPub *MessagePublisherImpl) PublishMessage(ctx context.Context, msg *Message) error {
	return msgPub.publisher.Publish(MessagePubTopic, message.NewMessage(
		watermill.NewUUID(),
		msg.Encode(),
	))
}
