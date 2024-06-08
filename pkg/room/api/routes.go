package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/omran95/chat-app/pkg/common"
	"github.com/omran95/chat-app/pkg/config"
	"github.com/omran95/chat-app/pkg/room"
	"gopkg.in/olahol/melody.v1"
)

var WsConn MelodyConn

type MelodyConn struct {
	*melody.Melody
}

func NewWebSocketConnection() MelodyConn {
	melody := melody.New()
	WsConn = MelodyConn{melody}
	return WsConn
}

type HttpServer struct {
	port        string
	name        string
	httpServer  *http.Server
	wsCon       MelodyConn
	engine      *gin.Engine
	logger      common.HttpLog
	roomService room.RoomService
}

func NewGinEngine(logger common.HttpLog, config config.Config) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(common.CorsMiddleware())
	engine.Use(common.LoggingMiddleware(logger))
	engine.Use(common.MaxConnectionsAllowed(config.Room.Http.Server.MaxConn))
	return engine
}

func NewHttpServer(name string, logger common.HttpLog, engine *gin.Engine, ws MelodyConn, config config.Config) *HttpServer {
	return &HttpServer{
		name:   name,
		logger: logger,
		engine: engine,
		wsCon:  ws,
		port:   config.Room.Http.Server.Port,
	}
}

func (server *HttpServer) RegisterRoutes() {
	roomGroup := server.engine.Group("/api/room")
	{
		roomGroup.POST("", server.CreateRoom)
	}
}

func (server *HttpServer) Run() {
	go func() {
		addr := ":" + server.port
		server.httpServer = &http.Server{
			Addr: addr,
		}
		server.logger.Info("Room HTTP server listening", slog.String("addr", addr))
		err := server.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			server.logger.Error(err.Error())
			os.Exit(1)
		}
	}()
}

func (server *HttpServer) GracefulStop(ctx context.Context) error {
	err := WsConn.Close()
	if err != nil {
		return err
	}
	err = server.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
