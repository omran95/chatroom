package room

import "encoding/json"

type CreateRoomDTO struct {
	Name      string `json:"name" binding:"required"`
	Protected bool   `json:"protected"`
	Password  string `json:"password"`
}

func (dto *CreateRoomDTO) isValid() bool {
	if dto.Protected && dto.Password == "" {
		return false
	}
	return true
}

type RoomID = uint64

type RoomPresenter struct {
	ID        RoomID `json:"room_id"`
	Name      string `json:"name"`
	Protected bool   `json:"protected"`
}

type Room struct {
	ID        RoomID `json:"room_id"`
	Name      string `json:"name"`
	Protected bool   `json:"protected"`
	Password  string `json:"password"`
}

func (room *Room) FromDTO(dto CreateRoomDTO) {
	room.Name = dto.Name
	room.Protected = dto.Protected
	room.Password = dto.Password
}

func (room *Room) ToPresenter() *RoomPresenter {
	return &RoomPresenter{
		ID:        room.ID,
		Name:      room.Name,
		Protected: room.Protected,
	}
}

type Action string

var (
	JoinedMessage    Action = "joined"
	IsTypingMessage  Action = "istyping"
	EndTypingMessage Action = "endtyping"
	LeftMessage      Action = "left"
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
