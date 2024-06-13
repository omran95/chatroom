package subscriber

import (
	"context"
	"strconv"

	"github.com/omran95/chatroom/pkg/infrastructure"
)

var redisPrefix = "subscriber"

type SubscriberRepo interface {
	AddRoomSubscriber(ctx context.Context, roomId uint64, userName, subscriber string) error
	RemoveRoomSubscriber(ctx context.Context, roomId uint64, userName string) error
	GetRoomSubscribers(ctx context.Context, roomId uint64) (map[string]struct{}, error)
}

type SubscriberRepoImpl struct {
	cache infrastructure.RedisCache
}

func NewSubscriberRepo(cache infrastructure.RedisCache) *SubscriberRepoImpl {
	return &SubscriberRepoImpl{cache}
}

func (repo *SubscriberRepoImpl) AddRoomSubscriber(ctx context.Context, roomID uint64, userName, subscriber string) error {
	key := constructRoomKey(roomID)
	return repo.cache.HSet(ctx, key, userName, subscriber)
}

func (repo *SubscriberRepoImpl) RemoveRoomSubscriber(ctx context.Context, roomID uint64, userName string) error {
	key := constructRoomKey(roomID)
	return repo.cache.HDel(ctx, key, userName)
}

func (repo *SubscriberRepoImpl) GetRoomSubscribers(ctx context.Context, roomID uint64) (map[string]struct{}, error) {
	key := constructRoomKey(roomID)

	roomSubscribers, err := repo.cache.HGetAll(ctx, key)
	if err != nil {
		return nil, err
	}

	// Using Map for uniquness
	subscribers := map[string]struct{}{}

	for _, subscriber := range roomSubscribers {
		subscribers[subscriber] = struct{}{}
	}

	return subscribers, nil
}

func constructRoomKey(roomID uint64) string {
	return redisPrefix + ":" + strconv.FormatUint(roomID, 10)
}
