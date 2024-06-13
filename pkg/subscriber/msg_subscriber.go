package subscriber

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/omran95/chatroom/pkg/room"
)

type MessageSubscriber struct {
	router            *message.Router
	subscriber        message.Subscriber
	subscriberService SubscriberService
}

func NewMessageSubscriber(router *message.Router, subscriber message.Subscriber, subscriberService SubscriberService) (*MessageSubscriber, error) {
	return &MessageSubscriber{
		router:            router,
		subscriber:        subscriber,
		subscriberService: subscriberService,
	}, nil
}

func (subscriber *MessageSubscriber) HandleIncomingMessage(msg *message.Message) error {
	message, err := DecodeToMessage(msg.Payload)
	if err != nil {
		return err
	}
	return subscriber.subscriberService.NotifySubscribers(msg.Context(), *message)
}

func (subscriber *MessageSubscriber) RegisterHandler() {
	subscriber.router.AddNoPublisherHandler(
		"subscriber_message_handler",
		room.MessagePubTopic,
		subscriber.subscriber,
		subscriber.HandleIncomingMessage,
	)
}

func (subscriber *MessageSubscriber) Run() error {
	return subscriber.router.Run(context.Background())
}

func (subscriber *MessageSubscriber) GracefulStop() error {
	return subscriber.router.Close()
}

func DecodeToMessage(data []byte) (*room.Message, error) {
	var msg room.Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
