package room

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omran95/chatroom/pkg/common"
	"github.com/omran95/chatroom/pkg/config"
	"github.com/redis/go-redis/v9"
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
	port                  string
	name                  string
	httpServer            *http.Server
	wsCon                 MelodyConn
	engine                *gin.Engine
	logger                common.HttpLog
	roomService           RoomService
	msgSubscriber         *MessageSubscriber
	rateLimiterMiddleware *RateLimiterMiddleware
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

func NewHttpServer(name string, logger common.HttpLog, engine *gin.Engine, ws MelodyConn, config *config.Config, roomService RoomService, msgSubscriber *MessageSubscriber, redisClient redis.UniversalClient) (*HttpServer, error) {
	// FillingRatePerRequest (RPS), bucketSize, expiration
	createRoomsrateLimiter, err := common.NewRateLimiter(redisClient, 1, 30, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("error creating room rate limiter: %w", err)
	}
	rateLimiterMiddleware := NewRateLimiterMiddleware(*createRoomsrateLimiter)
	return &HttpServer{
		name:                  name,
		logger:                logger,
		engine:                engine,
		wsCon:                 ws,
		port:                  config.Room.Http.Server.Port,
		roomService:           roomService,
		msgSubscriber:         msgSubscriber,
		rateLimiterMiddleware: rateLimiterMiddleware,
	}, nil
}

func (server *HttpServer) RegisterRoutes() {
	server.msgSubscriber.RegisterHandler()
	roomGroup := server.engine.Group("/api/rooms")
	{
		roomGroup.POST("", server.rateLimiterMiddleware.LimitCreateRooms, server.CreateRoom)
		roomGroup.GET("/:id", server.RequestToJoinRoom)
	}
	server.wsCon.HandleConnect(server.HandleRoomOnJoin)
	server.wsCon.HandleClose(server.HandleRoomOnLeave)
	server.wsCon.HandleMessage(server.HandleOnMessage)
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
	go func() {
		err := server.msgSubscriber.Run()
		if err != nil {
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
	err = server.msgSubscriber.GracefulStop()
	if err != nil {
		return err
	}
	return nil
}
