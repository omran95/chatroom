package subscriber

import "context"

type SubscriberService interface {
	AddRoomSubscriber(ctx context.Context, roomId uint64, userName, subscriber string) error
	RemoveRoomSubscriber(ctx context.Context, roomId uint64, userName string) error
}

type SubscriberServiceImpl struct {
	subscriberRepo SubscriberRepo
}

func NewSubscriberService(subscriberRepo SubscriberRepo) *SubscriberServiceImpl {
	return &SubscriberServiceImpl{subscriberRepo}
}

func (service *SubscriberServiceImpl) AddRoomSubscriber(ctx context.Context, roomId uint64, userName, subscriber string) error {
	return service.subscriberRepo.AddRoomSubscriber(ctx, roomId, userName, subscriber)
}

func (service *SubscriberServiceImpl) RemoveRoomSubscriber(ctx context.Context, roomId uint64, userName string) error {
	return service.subscriberRepo.RemoveRoomSubscriber(ctx, roomId, userName)
}
