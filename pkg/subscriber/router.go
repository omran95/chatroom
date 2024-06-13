package subscriber

import (
	"context"

	"github.com/omran95/chatroom/pkg/common"
)

type Router struct {
	grpcServer common.GrpcServer
}

func NewRouter(grpcServer common.GrpcServer) *Router {
	return &Router{grpcServer}
}

func (r *Router) Run() {
	r.grpcServer.Register()
	r.grpcServer.Run()

}
func (r *Router) GracefulStop(ctx context.Context) error {
	return r.grpcServer.GracefulStop()
}
