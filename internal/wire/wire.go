//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/omran95/chat-app/pkg/common"
	"github.com/omran95/chat-app/pkg/config"
	"github.com/omran95/chat-app/pkg/infrastructure"
	"github.com/omran95/chat-app/pkg/room"
	"github.com/omran95/chat-app/pkg/subscriber"
)

func InitializeRoomServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewHttpLog,
		common.NewSonyFlake,
		common.NewObservabilityInjector,
		infrastructure.NewCassandraSession,

		infrastructure.NewKafkaPublisher,
		room.NewSubscriberGrpcClient,

		room.NewMessagePublisher,
		wire.Bind(new(room.MessagePublisher), new(*room.MessagePublisherImpl)),

		room.NewRoomService,
		wire.Bind(new(room.RoomService), new(*room.RoomServiceImpl)),

		room.NewRoomRepo,
		wire.Bind(new(room.RoomRepo), new(*room.RoomRepoImpl)),

		room.NewWebSocketConnection,

		room.NewGinEngine,

		infrastructure.NewBrokerRouter,
		infrastructure.NewKafkaSubscriber,
		room.NewMessageSubscriber,

		room.NewHttpServer,
		wire.Bind(new(common.HttpServer), new(*room.HttpServer)),

		room.NewRouter,
		wire.Bind(new(common.Router), new(*room.Router)),
		common.NewServer,
	)
	return &common.Server{}, nil
}

func InitializeSubscriberServer(name string) (*common.Server, error) {
	wire.Build(
		config.NewConfig,
		common.NewGrpcLog,
		common.NewObservabilityInjector,

		infrastructure.NewRedisClient,
		infrastructure.NewRedisCacheImpl,
		wire.Bind(new(infrastructure.RedisCache), new(*infrastructure.RedisCacheImpl)),
		subscriber.NewSubscriberRepo,
		wire.Bind(new(subscriber.SubscriberRepo), new(*subscriber.SubscriberRepoImpl)),

		infrastructure.NewBrokerRouter,
		infrastructure.NewKafkaPublisher,
		infrastructure.NewKafkaSubscriber,

		subscriber.NewMessagePublisher,
		wire.Bind(new(subscriber.MessagePublisher), new(*subscriber.MessagePublisherImpl)),

		subscriber.NewMessageSubscriber,

		subscriber.NewSubscriberService,
		wire.Bind(new(subscriber.SubscriberService), new(*subscriber.SubscriberServiceImpl)),

		subscriber.NewGrpcServer,
		wire.Bind(new(common.GrpcServer), new(*subscriber.GrpcServer)),

		subscriber.NewRouter,
		wire.Bind(new(common.Router), new(*subscriber.Router)),
		common.NewServer,
	)
	return &common.Server{}, nil
}
