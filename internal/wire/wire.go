//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/omran95/chat-app/pkg/common"
	"github.com/omran95/chat-app/pkg/config"
	"github.com/omran95/chat-app/pkg/room"
)

func InitializeRoomServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewHttpLog,
		common.NewSonyFlake,

		room.NewRoomService,
		wire.Bind(new(room.RoomService), new(*room.RoomServiceImpl)),

		room.NewWebSocketConnection,

		room.NewGinEngine,

		room.NewHttpServer,
		wire.Bind(new(common.HttpServer), new(*room.HttpServer)),

		room.NewRouter,
		wire.Bind(new(common.Router), new(*room.Router)),
		common.NewServer,
	)
	return &common.Server{}, nil
}
