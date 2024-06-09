// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/omran95/chat-app/pkg/common"
	"github.com/omran95/chat-app/pkg/config"
	"github.com/omran95/chat-app/pkg/infrastructure"
	"github.com/omran95/chat-app/pkg/room"
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
	roomServiceImpl := room.NewRoomService(idGenerator, roomRepoImpl)
	httpServer := room.NewHttpServer(name, httpLog, engine, melodyConn, configConfig, roomServiceImpl)
	router := room.NewRouter(httpServer)
	observabilityInjector := common.NewObservabilityInjector(configConfig)
	server := common.NewServer(name, router, observabilityInjector)
	return server, nil
}
