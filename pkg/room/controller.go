package room

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/omran95/chat-app/pkg/common"
	"gopkg.in/olahol/melody.v1"
)

var sessRidKey = "sessRid"

func (server *HttpServer) CreateRoom(c *gin.Context) {
	room, err := server.roomService.CreateRoom(c)
	if err != nil {
		server.logger.Error(err.Error())
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}
	c.JSON(http.StatusCreated, room)
}

func (server *HttpServer) JoinRoom(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	userName := c.Query("userName")

	if err != nil || userName == "" {
		response(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}
	exist, err := server.roomService.RoomExist(c, roomID)

	if err != nil {
		server.logger.Error(err.Error())
		response(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	if !exist {
		response(c, http.StatusNotFound, common.ErrRoomNotFound)
		return
	}

	if err := server.wsCon.HandleRequest(c.Writer, c.Request); err != nil {
		server.logger.Error("upgrade websocket error: " + err.Error())
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}
}

func (server *HttpServer) HandleRoomOnJoin(wsSession *melody.Session) {
	roomID, userName := extractWsParams(wsSession)
	err := server.initializeChatSession(wsSession, roomID, userName)
	if err != nil {
		server.logger.Error(err.Error())
		return
	}

	if err := server.roomService.BroadcastConnectMessage(context.Background(), roomID, userName); err != nil {
		server.logger.Error(err.Error())
		return
	}
}

func (server *HttpServer) HandleRoomOnLeave(wsSession *melody.Session) {
	roomID, userName := extractWsParams(wsSession)
	err := server.roomService.RemoveRoomSubscriber(context.Background(), roomID, userName)
	if err != nil {
		server.logger.Error(err.Error())
		return
	}
	if err := server.roomService.BroadcastLeaveMessage(context.Background(), roomID, userName); err != nil {
		server.logger.Error(err.Error())
		return
	}
}

func (server *HttpServer) initializeChatSession(sess *melody.Session, roomID RoomID, userName string) error {
	ctx := context.Background()
	if err := server.roomService.AddRoomSubscriber(ctx, roomID, userName, server.msgSubscriber.topic); err != nil {
		return err
	}
	sess.Set(sessRidKey, roomID)
	return nil
}

func extractWsParams(wsSession *melody.Session) (roomID RoomID, userName string) {
	userName = wsSession.Request.URL.Query().Get("userName")
	// path e.g. /api/rooms/:roomID
	pathParts := strings.Split(wsSession.Request.URL.Path, "/")
	roomID, _ = strconv.ParseUint(pathParts[len(pathParts)-1], 10, 64)
	return
}

func response(c *gin.Context, httpCode int, err error) {
	message := err.Error()
	c.JSON(httpCode, common.ErrResponse{
		Message: message,
	})
}
