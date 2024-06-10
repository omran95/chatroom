package room

import (
	"context"
	"database/sql"

	"github.com/gocql/gocql"
)

type RoomRepo interface {
	CreateRoom(ctx context.Context, room Room) error
	RoomExist(ctx context.Context, roomID RoomID) (bool, error)
	IsProtected(ctx context.Context, roomID RoomID) (bool, error)
	IsValidPassword(ctx context.Context, roomID RoomID, password string) (bool, error)
}

type RoomRepoImpl struct {
	cassandraSession *gocql.Session
}

func NewRoomRepo(cassandraSession *gocql.Session) *RoomRepoImpl {
	return &RoomRepoImpl{cassandraSession}
}

func (repo *RoomRepoImpl) CreateRoom(ctx context.Context, room Room) error {
	query := "insert into rooms (id, name, protected, password) values (?, ?, ?, ?)"
	stmt := repo.cassandraSession.Query(query, room.ID, room.Name, room.Protected, room.Password).WithContext(ctx)
	if err := stmt.Exec(); err != nil {
		return err
	}
	return nil
}

func (repo *RoomRepoImpl) RoomExist(ctx context.Context, roomID RoomID) (bool, error) {
	var id RoomID
	err := repo.cassandraSession.Query("select id from rooms where id = ?", roomID).WithContext(ctx).Idempotent(true).Scan(&id)

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

func (repo *RoomRepoImpl) IsProtected(ctx context.Context, roomID RoomID) (bool, error) {
	var isProtected bool
	err := repo.cassandraSession.Query("select protected from rooms where id = ?", roomID).WithContext(ctx).Idempotent(true).Scan(&isProtected)
	if err != nil {
		return false, err
	}
	return isProtected, nil
}

func (repo *RoomRepoImpl) IsValidPassword(ctx context.Context, roomID RoomID, password string) (bool, error) {
	var roomPassword string
	err := repo.cassandraSession.Query("select password from rooms where id = ?", roomID).WithContext(ctx).Idempotent(true).Scan(&roomPassword)
	if err != nil {
		return false, err
	}
	return roomPassword == password, nil
}
