package subscriber

import (
	"context"
	"strconv"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/omran95/chat-app/pkg/infrastructure"
)

var redisPrefix = "subscriber"

type SubscriberRepo interface {
	AddRoomSubscriber(ctx context.Context, roomId uint64, userName, subscriber string) error
	RemoveRoomSubscriber(ctx context.Context, roomId uint64, userName string) error
}

type SubscriberRepoImpl struct {
	cache            infrastructure.RedisCache
	messagePublisher message.Publisher
}

func NewSubscriberRepo(cache infrastructure.RedisCache, messagePublisher message.Publisher) *SubscriberRepoImpl {
	return &SubscriberRepoImpl{cache, messagePublisher}
}

func (repo *SubscriberRepoImpl) AddRoomSubscriber(ctx context.Context, roomID uint64, userName, subscriber string) error {
	key := constructRoomKey(roomID)
	return repo.cache.HSet(ctx, key, userName, subscriber)
}

func (repo *SubscriberRepoImpl) RemoveRoomSubscriber(ctx context.Context, roomID uint64, userName string) error {
	key := constructRoomKey(roomID)
	return repo.cache.HDel(ctx, key, userName)
}

func constructRoomKey(roomID uint64) string {
	return redisPrefix + ":" + strconv.FormatUint(roomID, 10)
}
