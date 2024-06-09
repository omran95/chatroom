package room

import (
	"context"

	"github.com/gocql/gocql"
)

type RoomRepo interface {
	CreateRoom(ctx context.Context, roomID RoomID) error
}

type RoomRepoImpl struct {
	cassandraSession *gocql.Session
}

func NewRoomRepo(cassandraSession *gocql.Session) *RoomRepoImpl {
	return &RoomRepoImpl{cassandraSession}
}

func (repo *RoomRepoImpl) CreateRoom(ctx context.Context, roomID RoomID) error {
	if err := repo.cassandraSession.Query("insert into rooms (id) values (?)", roomID).WithContext(ctx).Exec(); err != nil {
		return err
	}
	return nil
}
