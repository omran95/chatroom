package room

import "encoding/json"

type RoomID = uint64

type Room struct {
	ID RoomID `json:"room_id"`
}

type Action string

var (
	JoinedMessage    Action = "joined"
	IsTypingMessage  Action = "istyping"
	EndTypingMessage Action = "endtyping"
	LeavedMessage    Action = "leaved"
)

type MessageID = uint64

const (
	EventText = iota
	EventAction
	EventSeen
	EventFile
)

type Message struct {
	ID       MessageID `json:"message_id"`
	Event    int       `json:"event"`
	RoomID   RoomID    `json:"room_id"`
	UserName string    `json:"username"`
	Payload  string    `json:"payload"`
	Seen     bool      `json:"seen"`
	Time     int64     `json:"time"`
}

func (m *Message) Encode() []byte {
	result, _ := json.Marshal(m)
	return result
}
