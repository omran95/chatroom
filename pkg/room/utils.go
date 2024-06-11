package room

import "encoding/json"

func decodeToMessage(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func decodeToRoomAuth(data []byte) (*RoomAuth, error) {
	var auth RoomAuth
	if err := json.Unmarshal(data, &auth); err != nil {
		return nil, err
	}
	return &auth, nil
}
