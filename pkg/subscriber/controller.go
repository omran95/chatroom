package subscriber

import (
	"context"

	subscriberpb "github.com/omran95/chat-app/pkg/subscriber/proto"
)

func (grpc *GrpcServer) AddRoomSubscriber(ctx context.Context, req *subscriberpb.AddRoomSubscriberRequest) (*subscriberpb.AddRoomSubscriberResponse, error) {

	return &subscriberpb.AddRoomSubscriberResponse{}, nil
}

func (grpc *GrpcServer) RemoveRoomSubscriber(ctx context.Context, req *subscriberpb.RemoveRoomSubscriberRequest) (*subscriberpb.RemoveRoomSubscriberResponse, error) {

	return &subscriberpb.RemoveRoomSubscriberResponse{}, nil
}
