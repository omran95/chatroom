package room

import (
	"context"
	"strconv"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

var MessagePubTopic = "chat.msg.pub"

type MessagePublisher interface {
	PublishMessage(ctx context.Context, message Message) error
}

type MessagePublisherImpl struct {
	publisher message.Publisher
}

func NewMessagePublisher(publisher message.Publisher) *MessagePublisherImpl {
	return &MessagePublisherImpl{publisher}
}

func (msgPub *MessagePublisherImpl) PublishMessage(ctx context.Context, msg Message) error {
	kafkaMessage := message.NewMessage(
		watermill.NewUUID(),
		msg.Encode(),
	)
	//partition kafka topic based on roomID
	kafkaMessage.Metadata.Set("partition_key", strconv.FormatUint(msg.RoomID, 10))
	return msgPub.publisher.Publish(MessagePubTopic, kafkaMessage)
}
