package subscriber

import (
	"context"

	subscriberpb "github.com/omran95/chat-app/pkg/subscriber/proto"
)

func (grpc *GrpcServer) AddRoomSubscriber(ctx context.Context, req *subscriberpb.AddRoomSubscriberRequest) (*subscriberpb.AddRoomSubscriberResponse, error) {
	err := grpc.subscriberService.AddRoomSubscriber(ctx, req.RoomId, req.Username, req.Subscriber)
	if err != nil {
		grpc.logger.Error(err.Error())
		return nil, err
	}
	return &subscriberpb.AddRoomSubscriberResponse{}, nil
}

func (grpc *GrpcServer) RemoveRoomSubscriber(ctx context.Context, req *subscriberpb.RemoveRoomSubscriberRequest) (*subscriberpb.RemoveRoomSubscriberResponse, error) {
	err := grpc.subscriberService.RemoveRoomSubscriber(ctx, req.RoomId, req.Username)
	if err != nil {
		grpc.logger.Error(err.Error())
		return nil, err
	}
	return &subscriberpb.RemoveRoomSubscriberResponse{}, nil
}
