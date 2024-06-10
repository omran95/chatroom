package room

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/omran95/chat-app/pkg/config"

	"gopkg.in/olahol/melody.v1"
)

type MessageSubscriber struct {
	topic      string
	router     *message.Router
	subscriber message.Subscriber
	ws         MelodyConn
}

func NewMessageSubscriber(router *message.Router, config *config.Config, subscriber message.Subscriber, ws MelodyConn) (*MessageSubscriber, error) {
	return &MessageSubscriber{
		topic:      config.Room.MessageSubscriber.Topic,
		router:     router,
		subscriber: subscriber,
		ws:         ws,
	}, nil
}

func (subscriber *MessageSubscriber) HandleIncomingMessage(msg *message.Message) error {
	message, err := decodeToMessage([]byte(msg.Payload))
	if err != nil {
		return err
	}
	return subscriber.broadcast(message)
}

func (subscriber *MessageSubscriber) RegisterHandler() {
	subscriber.router.AddNoPublisherHandler(
		"room_message_handler",
		subscriber.topic,
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

func (subscriber *MessageSubscriber) broadcast(message *Message) error {
	return subscriber.ws.BroadcastFilter(message.Encode(), func(sess *melody.Session) bool {
		roomID, exist := sess.Get(sessRidKey)
		if !exist {
			return false
		}
		return message.RoomID == (roomID.(uint64))
	})
}

func decodeToMessage(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
