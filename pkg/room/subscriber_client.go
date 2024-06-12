package room

import (
	"github.com/omran95/chat-app/pkg/config"
	"github.com/omran95/chat-app/pkg/infrastructure"
	"google.golang.org/grpc"
)

type SubscriberGrpcClient struct {
	Conn *grpc.ClientConn
}

func NewSubscriberGrpcClient(config *config.Config) (*SubscriberGrpcClient, error) {
	conn, err := infrastructure.InitializeGrpcClient(config.Room.Grpc.Client.Subscriber.Endpoint)
	if err != nil {
		return nil, err
	}
	SubscriberConn := &SubscriberGrpcClient{
		Conn: conn,
	}
	return SubscriberConn, nil
}
