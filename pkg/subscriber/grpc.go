package subscriber

import (
	"net"
	"os"

	"log/slog"

	"github.com/omran95/chatroom/pkg/common"
	"github.com/omran95/chatroom/pkg/config"
	"github.com/omran95/chatroom/pkg/infrastructure"

	subscriberpb "github.com/omran95/chatroom/pkg/subscriber/proto"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	port              string
	logger            common.GrpcLog
	server            *grpc.Server
	subscriberService SubscriberService
	msgSubscriber     *MessageSubscriber
	subscriberpb.UnimplementedSubscriberServiceServer
}

func NewGrpcServer(name string, logger common.GrpcLog, config *config.Config, subscriberService SubscriberService, msgSubscriber *MessageSubscriber) *GrpcServer {
	grpc := &GrpcServer{
		port:              config.Subscriber.Grpc.Server.Port,
		logger:            logger,
		subscriberService: subscriberService,
		msgSubscriber:     msgSubscriber,
	}
	grpc.server = infrastructure.InitializeGrpcServer(name, grpc.logger)
	return grpc
}

func (grpc *GrpcServer) Register() {
	grpc.msgSubscriber.RegisterHandler()
	subscriberpb.RegisterSubscriberServiceServer(grpc.server, grpc)
}

func (grpc *GrpcServer) Run() {
	go func() {
		addr := "0.0.0.0:" + grpc.port
		grpc.logger.Info("grpc server listening", slog.String("addr", addr))
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			grpc.logger.Error(err.Error())
			os.Exit(1)
		}
		if err := grpc.server.Serve(lis); err != nil {
			grpc.logger.Error(err.Error())
			os.Exit(1)
		}
	}()
	go func() {
		err := grpc.msgSubscriber.Run()
		if err != nil {
			grpc.logger.Error(err.Error())
			os.Exit(1)
		}
	}()
}

func (grpc *GrpcServer) GracefulStop() error {
	grpc.server.GracefulStop()
	return grpc.msgSubscriber.GracefulStop()
}
