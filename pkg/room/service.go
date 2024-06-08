package room

import "context"

type RoomID = string

type RoomService interface {
	CreateRoom(ctx context.Context) (RoomID, error)
}
