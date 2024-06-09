package room

import (
	"context"
	"fmt"
	"time"

	"github.com/omran95/chat-app/pkg/common"
)

type RoomService interface {
	CreateRoom(ctx context.Context) (Room, error)
	RoomExist(ctx context.Context, roomID RoomID) (bool, error)
	BroadcastConnectMessage(ctx context.Context, roomID RoomID, userName string) error
}

type RoomServiceImpl struct {
	snowFlake        common.IDGenerator
	roomRepo         RoomRepo
	messagePublisher MessagePublisher
}

func NewRoomService(snowflake common.IDGenerator, roomRepo RoomRepo, messagePublisher MessagePublisher) *RoomServiceImpl {
	return &RoomServiceImpl{snowflake, roomRepo, messagePublisher}
}

func (service *RoomServiceImpl) CreateRoom(ctx context.Context) (Room, error) {
	roomID, err := service.snowFlake.NextID()
	if err != nil {
		return Room{}, fmt.Errorf("error create snowflake ID for new room: %w", err)
	}
	if err := service.roomRepo.CreateRoom(ctx, roomID); err != nil {
		return Room{}, fmt.Errorf("error creating room: %w", err)
	}
	return Room{ID: roomID}, nil
}

func (service *RoomServiceImpl) RoomExist(ctx context.Context, roomID RoomID) (bool, error) {
	room, err := service.roomRepo.RoomExist(ctx, roomID)
	if err != nil {
		return false, fmt.Errorf("error checking room existence: %w", err)
	}
	return room, nil
}

func (service *RoomServiceImpl) BroadcastConnectMessage(ctx context.Context, roomID RoomID, userName string) error {

	return service.BroadcastActionMessage(ctx, roomID, userName, JoinedMessage)
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
