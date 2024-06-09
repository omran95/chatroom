//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/omran95/chat-app/pkg/common"
	"github.com/omran95/chat-app/pkg/config"
	"github.com/omran95/chat-app/pkg/infrastructure"
	"github.com/omran95/chat-app/pkg/room"
)

func InitializeRoomServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewHttpLog,
		common.NewSonyFlake,
		common.NewObservabilityInjector,
		infrastructure.NewCassandraSession,

		infrastructure.NewKafkaPublisher,

		room.NewMessagePublisher,
		wire.Bind(new(room.MessagePublisher), new(*room.MessagePublisherImpl)),

		room.NewRoomService,
		wire.Bind(new(room.RoomService), new(*room.RoomServiceImpl)),

		room.NewRoomRepo,
		wire.Bind(new(room.RoomRepo), new(*room.RoomRepoImpl)),

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
