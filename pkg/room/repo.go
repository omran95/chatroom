package room

import (
	"context"
	"database/sql"

	"github.com/gocql/gocql"
)

type RoomRepo interface {
	CreateRoom(ctx context.Context, roomID RoomID) error
	RoomExist(ctx context.Context, roomID RoomID) (bool, error)
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

func (repo *RoomRepoImpl) RoomExist(ctx context.Context, roomID RoomID) (bool, error) {
	var id RoomID
	err := repo.cassandraSession.Query("select * from rooms where id = ?", roomID).WithContext(ctx).Idempotent(true).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			// Room does not exist
			return false, nil
		}
		// Return error for other errors
		return false, err
	}

	return true, nil
}
