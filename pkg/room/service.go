package room

import (
	"context"
	"fmt"

	"github.com/omran95/chat-app/pkg/common"
)

type RoomService interface {
	CreateRoom(ctx context.Context) (Room, error)
}

type RoomServiceImpl struct {
	snowFlake common.IDGenerator
	roomRepo  RoomRepo
}

func NewRoomService(snowflake common.IDGenerator, roomRepo RoomRepo) *RoomServiceImpl {
	return &RoomServiceImpl{snowflake, roomRepo}
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
