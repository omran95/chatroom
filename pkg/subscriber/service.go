package subscriber

import (
	"context"

	"github.com/omran95/chatroom/pkg/room"
)

type SubscriberService interface {
	AddRoomSubscriber(ctx context.Context, roomId uint64, userName, subscriber string) error
	RemoveRoomSubscriber(ctx context.Context, roomId uint64, userName string) error
	NotifySubscribers(ctx context.Context, message room.Message) error
}

type SubscriberServiceImpl struct {
	msgPublisher   MessagePublisher
	subscriberRepo SubscriberRepo
}

func NewSubscriberService(subscriberRepo SubscriberRepo, msgPublisher MessagePublisher) *SubscriberServiceImpl {
	return &SubscriberServiceImpl{msgPublisher, subscriberRepo}
}

func (service *SubscriberServiceImpl) AddRoomSubscriber(ctx context.Context, roomId uint64, userName, subscriber string) error {
	return service.subscriberRepo.AddRoomSubscriber(ctx, roomId, userName, subscriber)
}

func (service *SubscriberServiceImpl) RemoveRoomSubscriber(ctx context.Context, roomId uint64, userName string) error {
	return service.subscriberRepo.RemoveRoomSubscriber(ctx, roomId, userName)
}

func (service *SubscriberServiceImpl) NotifySubscribers(ctx context.Context, message room.Message) error {
	roomSubscribers, err := service.subscriberRepo.GetRoomSubscribers(ctx, message.RoomID)
	if err != nil {
		return err
	}
	return service.msgPublisher.PublishToSubscribers(ctx, roomSubscribers, message)
}
