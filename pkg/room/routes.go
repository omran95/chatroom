package room

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/omran95/chat-app/pkg/common"
	"github.com/omran95/chat-app/pkg/config"
	"gopkg.in/olahol/melody.v1"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	prommiddleware "github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
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
	roomService RoomService
}

func NewGinEngine(name string, logger common.HttpLog, config *config.Config) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(common.CorsMiddleware())
	engine.Use(common.LoggingMiddleware(logger))
	engine.Use(common.MaxConnectionsAllowed(config.Room.Http.Server.MaxConn))
	mdlw := prommiddleware.New(prommiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{
			Prefix: name,
		}),
	})
	engine.Use(ginmiddleware.Handler("", mdlw))
	return engine
}

func NewHttpServer(name string, logger common.HttpLog, engine *gin.Engine, ws MelodyConn, config *config.Config, roomService RoomService) *HttpServer {
	return &HttpServer{
		name:        name,
		logger:      logger,
		engine:      engine,
		wsCon:       ws,
		port:        config.Room.Http.Server.Port,
		roomService: roomService,
	}
}

func (server *HttpServer) RegisterRoutes() {
	roomGroup := server.engine.Group("/api/rooms")
	{
		roomGroup.POST("", server.CreateRoom)
	}
}

func (server *HttpServer) Run() {
	go func() {
		addr := ":" + server.port
		server.httpServer = &http.Server{
			Addr:    addr,
			Handler: common.NewOtelHttpHandler(server.engine, server.name+"_http"),
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
