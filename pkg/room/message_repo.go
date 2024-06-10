package room

import (
	"context"

	"github.com/gocql/gocql"
)

type MessageRepo interface {
	InesrtMessage(ctx context.Context, msg Message) error
}

type MessageRepoImpl struct {
	cassandraSession *gocql.Session
	insertStmt       *gocql.Query
}

func NewMessageRepo(cassandraSession *gocql.Session) *MessageRepoImpl {
	query := "insert into messages (id, event, room_id, username, payload, seen, timestamp) values (?, ?, ?, ?, ?, ?, ?)"
	preparedStmt := cassandraSession.Query(query)
	return &MessageRepoImpl{cassandraSession, preparedStmt}
}

func (msgRepo *MessageRepoImpl) InesrtMessage(ctx context.Context, msg Message) error {
	stmt := msgRepo.insertStmt.Bind(msg.ID, msg.Event, msg.RoomID, msg.UserName, msg.Payload, msg.Seen, msg.Time).WithContext(ctx)

	if err := stmt.Exec(); err != nil {
		return err
	}
	return nil
}
