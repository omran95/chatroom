package common

import (
	"errors"

	"github.com/sony/sonyflake"
)

// IDGenerator is the inteface for generatring unique ID
type IDGenerator interface {
	NextID() (uint64, error)
}

func NewSonyFlake() (IDGenerator, error) {
	var settings sonyflake.Settings
	snowFlake := sonyflake.NewSonyflake(settings)
	if snowFlake == nil {
		return nil, errors.New("sonyflake not created")
	}
	return snowFlake, nil
}
