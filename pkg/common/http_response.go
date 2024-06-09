package common

import (
	"errors"
)

var (
	ErrInvalidParam = errors.New("invalid parameter")
	ErrServer       = errors.New("internal server error")
	ErrRoomNotFound = errors.New("room not found")
)

// ErrResponse is the error response type
type ErrResponse struct {
	Message string `json:"msg"`
}
