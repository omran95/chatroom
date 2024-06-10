package room

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/omran95/chat-app/pkg/common"
	"github.com/omran95/chat-app/pkg/infrastructure"
	subscriberpb "github.com/omran95/chat-app/pkg/subscriber/proto"
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
}

type RoomServiceImpl struct {
	snowFlake                 common.IDGenerator
	roomRepo                  RoomRepo
	messagePublisher          MessagePublisher
	AddRoomSubscriberEndpoint endpoint.Endpoint
	RemoveSubscriberEndpoint  endpoint.Endpoint
}

func NewRoomService(snowflake common.IDGenerator, roomRepo RoomRepo, messagePublisher MessagePublisher, subscriberClient *SubscriberGrpcClient) *RoomServiceImpl {
	AddRoomSubscriberEndpoint := infrastructure.NewGrpcEndpoint(subscriberClient.Conn, "subscriber", "proto.SubscriberService", "AddRoomSubscriber", &subscriberpb.AddRoomSubscriberResponse{})
	RemoveRoomSubscriberEndpoint := infrastructure.NewGrpcEndpoint(subscriberClient.Conn, "subscriber", "proto.SubscriberService", "RemoveRoomSubscriber", &subscriberpb.RemoveRoomSubscriberResponse{})

	return &RoomServiceImpl{snowflake, roomRepo, messagePublisher, AddRoomSubscriberEndpoint, RemoveRoomSubscriberEndpoint}
}

func (service *RoomServiceImpl) CreateRoom(ctx context.Context, dto CreateRoomDTO) (*RoomPresenter, error) {
	roomID, err := service.snowFlake.NextID()
	if err != nil {
		return nil, fmt.Errorf("error create snowflake ID for new room: %w", err)
	}
	room := &Room{ID: roomID}
	room.FromDTO(dto)
	if err := service.roomRepo.CreateRoom(ctx, *room); err != nil {
		return nil, fmt.Errorf("error creating room: %w", err)
	}
	return room.ToPresenter(), nil
}

func (service *RoomServiceImpl) RoomExist(ctx context.Context, roomID RoomID) (bool, error) {
	room, err := service.roomRepo.RoomExist(ctx, roomID)
	if err != nil {
		return false, fmt.Errorf("error checking room existence: %w", err)
	}
	return room, nil
}

func (service *RoomServiceImpl) IsRoomProtected(ctx context.Context, roomID RoomID) (bool, error) {
	protected, err := service.roomRepo.IsProtected(ctx, roomID)
	if err != nil {
		return false, fmt.Errorf("error checking room protection: %w", err)
	}
	return protected, nil
}

func (service *RoomServiceImpl) IsValidPassword(ctx context.Context, roomID RoomID, password string) (bool, error) {
	valid, err := service.roomRepo.IsValidPassword(ctx, roomID, password)
	if err != nil {
		return false, fmt.Errorf("error validating room passowrd: %w", err)
	}
	return valid, nil
}

func (service *RoomServiceImpl) BroadcastConnectMessage(ctx context.Context, roomID RoomID, userName string) error {
	return service.BroadcastActionMessage(ctx, roomID, userName, JoinedMessage)
}

func (service *RoomServiceImpl) BroadcastLeaveMessage(ctx context.Context, roomID RoomID, userName string) error {
	return service.BroadcastActionMessage(ctx, roomID, userName, LeavedMessage)
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
	if err := service.messagePublisher.PublishMessage(ctx, &msg); err != nil {
		return fmt.Errorf("error broadcast action message: %w", err)
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
