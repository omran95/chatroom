// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/omran95/chatroom/pkg/common"
	"github.com/omran95/chatroom/pkg/config"
	"github.com/omran95/chatroom/pkg/infrastructure"
	"github.com/omran95/chatroom/pkg/room"
	"github.com/omran95/chatroom/pkg/subscriber"
)

// Injectors from wire.go:

func InitializeRoomServer(name string) (*common.Server, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	httpLog, err := common.NewHttpLog(configConfig)
	if err != nil {
		return nil, err
	}
	engine := room.NewGinEngine(name, httpLog, configConfig)
	melodyConn := room.NewWebSocketConnection()
	idGenerator, err := common.NewSonyFlake()
	if err != nil {
		return nil, err
	}
	session, err := infrastructure.NewCassandraSession(configConfig)
	if err != nil {
		return nil, err
	}
	roomRepoImpl := room.NewRoomRepo(session)
	publisher, err := infrastructure.NewKafkaPublisher(configConfig)
	if err != nil {
		return nil, err
	}
	messagePublisherImpl := room.NewMessagePublisher(publisher)
	subscriberGrpcClient, err := room.NewSubscriberGrpcClient(configConfig)
	if err != nil {
		return nil, err
	}
	messageRepoImpl := room.NewMessageRepo(session)
	roomServiceImpl := room.NewRoomService(idGenerator, roomRepoImpl, messagePublisherImpl, subscriberGrpcClient, messageRepoImpl)
	router, err := infrastructure.NewBrokerRouter(name)
	if err != nil {
		return nil, err
	}
	subscriber, err := infrastructure.NewKafkaSubscriber(configConfig)
	if err != nil {
		return nil, err
	}
	messageSubscriber, err := room.NewMessageSubscriber(router, configConfig, subscriber, melodyConn)
	if err != nil {
		return nil, err
	}
	universalClient, err := infrastructure.NewRedisClient(configConfig)
	if err != nil {
		return nil, err
	}
	httpServer, err := room.NewHttpServer(name, httpLog, engine, melodyConn, configConfig, roomServiceImpl, messageSubscriber, universalClient)
	if err != nil {
		return nil, err
	}
	roomRouter := room.NewRouter(httpServer)
	observabilityInjector := common.NewObservabilityInjector(configConfig)
	server := common.NewServer(name, roomRouter, observabilityInjector)
	return server, nil
}

func InitializeSubscriberServer(name string) (*common.Server, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	grpcLog, err := common.NewGrpcLog(configConfig)
	if err != nil {
		return nil, err
	}
	universalClient, err := infrastructure.NewRedisClient(configConfig)
	if err != nil {
		return nil, err
	}
	redisCacheImpl := infrastructure.NewRedisCacheImpl(universalClient)
	subscriberRepoImpl := subscriber.NewSubscriberRepo(redisCacheImpl)
	publisher, err := infrastructure.NewKafkaPublisher(configConfig)
	if err != nil {
		return nil, err
	}
	messagePublisherImpl := subscriber.NewMessagePublisher(publisher)
	subscriberServiceImpl := subscriber.NewSubscriberService(subscriberRepoImpl, messagePublisherImpl)
	router, err := infrastructure.NewBrokerRouter(name)
	if err != nil {
		return nil, err
	}
	messageSubscriber, err := infrastructure.NewKafkaSubscriber(configConfig)
	if err != nil {
		return nil, err
	}
	subscriberMessageSubscriber, err := subscriber.NewMessageSubscriber(router, messageSubscriber, subscriberServiceImpl)
	if err != nil {
		return nil, err
	}
	grpcServer := subscriber.NewGrpcServer(name, grpcLog, configConfig, subscriberServiceImpl, subscriberMessageSubscriber)
	subscriberRouter := subscriber.NewRouter(grpcServer)
	observabilityInjector := common.NewObservabilityInjector(configConfig)
	server := common.NewServer(name, subscriberRouter, observabilityInjector)
	return server, nil
}
