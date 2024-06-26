package room

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/omran95/chatroom/pkg/common"
	"github.com/omran95/chatroom/pkg/infrastructure"
	subscriberpb "github.com/omran95/chatroom/pkg/subscriber/proto"
	"golang.org/x/crypto/bcrypt"
)

type RoomService interface {
	CreateRoom(ctx context.Context, dto CreateRoomDTO) (*RoomPresenter, error)
	RoomExist(ctx context.Context, roomID RoomID) (bool, error)
	IsRoomProtected(ctx context.Context, roomID RoomID) (bool, error)
	IsValidPassword(ctx context.Context, roomID RoomID, password string) (bool, error)
	BroadcastConnectMessage(ctx context.Context, roomID RoomID, userName string) error
	BroadcastLeaveMessage(ctx context.Context, roomID RoomID, userName string) error
	AddRoomSubscriber(ctx context.Context, roomID RoomID, userName string, subscriberTopic string) error
	RemoveRoomSubscriber(ctx context.Context, roomID RoomID, userName string) error
	HandleNewMessage(ctx context.Context, msg Message) error
}

type RoomServiceImpl struct {
	snowFlake                 common.IDGenerator
	roomRepo                  RoomRepo
	messagePublisher          MessagePublisher
	AddRoomSubscriberEndpoint endpoint.Endpoint
	RemoveSubscriberEndpoint  endpoint.Endpoint
	messageRepo               MessageRepo
}

func NewRoomService(snowflake common.IDGenerator, roomRepo RoomRepo, messagePublisher MessagePublisher, subscriberClient *SubscriberGrpcClient, messageRepo MessageRepo) *RoomServiceImpl {
	AddRoomSubscriberEndpoint := infrastructure.NewGrpcEndpoint(subscriberClient.Conn, "subscriber", "proto.SubscriberService", "AddRoomSubscriber", &subscriberpb.AddRoomSubscriberResponse{})
	RemoveRoomSubscriberEndpoint := infrastructure.NewGrpcEndpoint(subscriberClient.Conn, "subscriber", "proto.SubscriberService", "RemoveRoomSubscriber", &subscriberpb.RemoveRoomSubscriberResponse{})
	return &RoomServiceImpl{snowflake, roomRepo, messagePublisher, AddRoomSubscriberEndpoint, RemoveRoomSubscriberEndpoint, messageRepo}
}

func (service *RoomServiceImpl) CreateRoom(ctx context.Context, dto CreateRoomDTO) (*RoomPresenter, error) {
	roomID, err := service.snowFlake.NextID()
	if err != nil {
		return nil, fmt.Errorf("error create snowflake ID for new room: %w", err)
	}
	room := &Room{ID: roomID}
	room.FromDTO(dto)
	if !room.Protected && room.Password != "" {
		room.Password = ""
	}
	if room.Protected {
		hashedPassword, err := hashPassword(room.Password)
		if err != nil {
			return nil, fmt.Errorf("error creating room: %w", err)
		}
		room.Password = hashedPassword
	}
	if err := service.roomRepo.CreateRoom(ctx, *room); err != nil {
		return nil, fmt.Errorf("error creating room: %w", err)
	}
	return room.ToPresenter(), nil
}

func (service *RoomServiceImpl) RoomExist(ctx context.Context, roomID RoomID) (bool, error) {
	exists, err := service.roomRepo.RoomExist(ctx, roomID)
	if err != nil {
		return false, fmt.Errorf("error checking room existence: %w", err)
	}
	return exists, nil
}

func (service *RoomServiceImpl) IsRoomProtected(ctx context.Context, roomID RoomID) (bool, error) {
	protected, err := service.roomRepo.IsProtected(ctx, roomID)
	if err != nil {
		return false, fmt.Errorf("error checking room protection: %w", err)
	}
	return protected, nil
}

func (service *RoomServiceImpl) IsValidPassword(ctx context.Context, roomID RoomID, password string) (bool, error) {
	roomPassword, err := service.roomRepo.GetRoomPassword(ctx, roomID)
	if err != nil {
		return false, fmt.Errorf("error validating passowrd: %w", err)
	}
	return validPassword(roomPassword, password), nil
}

func (service *RoomServiceImpl) BroadcastConnectMessage(ctx context.Context, roomID RoomID, userName string) error {
	return service.BroadcastActionMessage(ctx, roomID, userName, JoinedMessage)
}

func (service *RoomServiceImpl) BroadcastLeaveMessage(ctx context.Context, roomID RoomID, userName string) error {
	return service.BroadcastActionMessage(ctx, roomID, userName, LeftMessage)
}

func (service *RoomServiceImpl) BroadcastActionMessage(ctx context.Context, roomID RoomID, userName string, action Action) error {
	eventMessageID, err := service.snowFlake.NextID()
	if err != nil {
		return fmt.Errorf("error create snowflake ID for action message: %w", err)
	}
	msg := Message{
		ID:       eventMessageID,
		Event:    EventAction,
		RoomID:   roomID,
		UserName: userName,
		Payload:  string(action),
		Time:     time.Now().UnixMilli(),
	}
	if err := service.messagePublisher.PublishMessage(ctx, msg); err != nil {
		return fmt.Errorf("error broadcast action message: %w", err)
	}
	return nil
}

func (service *RoomServiceImpl) BroadcastTextMessage(ctx context.Context, roomID RoomID, userName string, payload string) error {
	messageID, err := service.snowFlake.NextID()
	if err != nil {
		return fmt.Errorf("error create snowflake ID for text message: %w", err)
	}
	msg := Message{
		ID:       messageID,
		Event:    EventText,
		RoomID:   roomID,
		UserName: userName,
		Payload:  payload,
		Time:     time.Now().UnixMilli(),
	}
	if err := service.messageRepo.InesrtMessage(ctx, msg); err != nil {
		return fmt.Errorf("error saving text message: %w", err)
	}
	if err := service.messagePublisher.PublishMessage(ctx, msg); err != nil {
		return fmt.Errorf("error broadcast text message: %w", err)
	}
	return nil
}

func (service *RoomServiceImpl) AddRoomSubscriber(ctx context.Context, roomID RoomID, userName string, subscriberTopic string) error {
	_, err := service.AddRoomSubscriberEndpoint(ctx, &subscriberpb.AddRoomSubscriberRequest{
		RoomId:          roomID,
		Username:        userName,
		SubscriberTopic: subscriberTopic,
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *RoomServiceImpl) RemoveRoomSubscriber(ctx context.Context, roomID RoomID, userName string) error {
	_, err := service.RemoveSubscriberEndpoint(ctx, &subscriberpb.RemoveRoomSubscriberRequest{
		RoomId:   roomID,
		Username: userName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *RoomServiceImpl) MarkSeen(ctx context.Context, roomID RoomID, userName string, seenMessageID MessageID) error {
	if err := service.messageRepo.MarkSeen(ctx, roomID, seenMessageID); err != nil {
		return err
	}
	messageID, err := service.snowFlake.NextID()
	if err != nil {
		return fmt.Errorf("error create snowflake ID for seen message: %w", err)
	}
	msg := Message{
		ID:       messageID,
		Event:    EventSeen,
		RoomID:   roomID,
		UserName: userName,
		Payload:  strconv.FormatUint(seenMessageID, 10),
		Seen:     true,
		Time:     time.Now().UnixMilli(),
	}
	if err := service.messagePublisher.PublishMessage(ctx, msg); err != nil {
		return fmt.Errorf("error broadcast seen message: %w", err)
	}
	return nil
}

func (service *RoomServiceImpl) HandleNewMessage(ctx context.Context, msg Message) error {
	switch msg.Event {
	case EventAction:
		return service.BroadcastActionMessage(ctx, msg.RoomID, msg.UserName, Action(msg.Payload))
	case EventText:
		return service.BroadcastTextMessage(ctx, msg.RoomID, msg.UserName, msg.Payload)
	case EventSeen:
		seenMessageID, err := strconv.ParseUint(msg.Payload, 10, 64)
		if err != nil {
			return err
		}
		return service.MarkSeen(ctx, msg.RoomID, msg.UserName, seenMessageID)
	}
	return nil
}

func hashPassword(password string) (string, error) {
	// Generate a bcrypt hash of the password with a cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func validPassword(roomPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(roomPassword), []byte(password))
	return err == nil
}
