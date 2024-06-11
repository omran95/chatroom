package room

import (
	"context"

	"github.com/gocql/gocql"
)

type MessageRepo interface {
	InesrtMessage(ctx context.Context, msg Message) error
	MarkSeen(ctx context.Context, RoomID RoomID, messageID MessageID) error
}

type MessageRepoImpl struct {
	cassandraSession *gocql.Session
	insertStmt       *gocql.Query
	markSeenStmt     *gocql.Query
}

func NewMessageRepo(cassandraSession *gocql.Session) *MessageRepoImpl {
	insertQuery := "insert into messages (id, event, room_id, username, payload, seen, timestamp) values (?, ?, ?, ?, ?, ?, ?)"
	preparedInsrtStmt := cassandraSession.Query(insertQuery)

	markSeenQuery := "UPDATE messages SET seen = ? WHERE room_id = ? AND id = ?"
	preparedmarkSeenStmt := cassandraSession.Query(markSeenQuery)

	return &MessageRepoImpl{cassandraSession, preparedInsrtStmt, preparedmarkSeenStmt}
}

func (msgRepo *MessageRepoImpl) InesrtMessage(ctx context.Context, msg Message) error {
	stmt := msgRepo.insertStmt.Bind(msg.ID, msg.Event, msg.RoomID, msg.UserName, msg.Payload, msg.Seen, msg.Time).WithContext(ctx)

	if err := stmt.Exec(); err != nil {
		return err
	}
	return nil
}

func (msgRepo *MessageRepoImpl) MarkSeen(ctx context.Context, RoomID RoomID, messageID MessageID) error {
	stmt := msgRepo.markSeenStmt.Bind(true, RoomID, messageID).Idempotent(true).WithContext(ctx)
	if err := stmt.Exec(); err != nil {
		return err
	}
	return nil
}
