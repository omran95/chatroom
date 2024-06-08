package room

import (
	"context"
	"fmt"

	"github.com/omran95/chat-app/pkg/common"
)

type RoomID = uint64

type RoomService interface {
	CreateRoom(ctx context.Context) (RoomID, error)
}

type RoomServiceImpl struct {
	snowFlake common.IDGenerator
}

func NewRoomService(snowflake common.IDGenerator) *RoomServiceImpl {
	return &RoomServiceImpl{snowflake}
}

func (service *RoomServiceImpl) CreateRoom(ctx context.Context) (RoomID, error) {
	roomID, err := service.snowFlake.NextID()
	if err != nil {
		return 0, fmt.Errorf("error create snowflake ID for new room: %w", err)
	}
	return roomID, nil
}
