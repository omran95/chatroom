package room

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/omran95/chatroom/pkg/common"
	"gopkg.in/olahol/melody.v1"
)

var sessRidKey = "sessRid"

func (server *HttpServer) CreateRoom(c *gin.Context) {
	var dto CreateRoomDTO
	if err := c.ShouldBindBodyWithJSON(&dto); err != nil {
		response(c, http.StatusBadRequest, err)
		return
	}

	if isValid := dto.isValid(); !isValid {
		response(c, http.StatusBadRequest, common.ErrInvalidParam)
		return
	}

	room, err := server.roomService.CreateRoom(c, dto)
	if err != nil {
		server.logger.Error(err.Error())
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}
	c.JSON(http.StatusCreated, room)
}

func (server *HttpServer) RequestToJoinRoom(c *gin.Context) {
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
	ctx := context.Background()
	roomID, userName := extractWsParams(wsSession)
	isProtectedRoom, err := server.roomService.IsRoomProtected(ctx, roomID)
	if err != nil {
		wsSession.CloseWithMsg(melody.FormatCloseMessage(500, "Error checking if the room is protected: "+err.Error()))
		return
	}
	if !isProtectedRoom {
		server.joinRoom(wsSession, roomID, userName)
		return
	}
	err = server.sendAuthRequiredMessage(wsSession)
	if err != nil {
		wsSession.CloseWithMsg(melody.FormatCloseMessage(500, "Error: "+err.Error()))
		return
	}
}

func (server *HttpServer) HandleRoomOnLeave(wsSession *melody.Session, n int, s string) error {
	roomID, userName := extractWsParams(wsSession)
	err := server.roomService.RemoveRoomSubscriber(context.Background(), roomID, userName)
	if err != nil {
		server.logger.Error(err.Error())
		return err
	}
	if err := server.roomService.BroadcastLeaveMessage(context.Background(), roomID, userName); err != nil {
		server.logger.Error(err.Error())
		return err
	}
	return nil
}

func (server *HttpServer) HandleOnMessage(wsSession *melody.Session, msg []byte) {
	roomID, userName := extractWsParams(wsSession)
	if authRequired := server.roomAuthRequired(wsSession); authRequired {
		server.AuthenticateRoom(wsSession, roomID, userName, msg)
		return
	}
	decodedMsg, err := decodeToMessage(msg)
	if err != nil {
		wsSession.CloseWithMsg(melody.FormatCloseMessage(400, "Invalid message"))
		return
	}
	decodedMsg.RoomID = roomID
	decodedMsg.UserName = userName
	err = server.roomService.HandleNewMessage(context.Background(), *decodedMsg)

	if err != nil {
		server.logger.Error(err.Error())
	}
}

func (server *HttpServer) roomAuthRequired(wsSession *melody.Session) bool {
	_, exists := wsSession.Get(sessRidKey)
	return !exists
}

func (server *HttpServer) AuthenticateRoom(wsSession *melody.Session, roomID RoomID, userName string, msg []byte) {
	password, err := extractPassword(msg)

	if err != nil || password == "" {
		wsSession.CloseWithMsg(melody.FormatCloseMessage(400, "Invalid password"))
		return

	}
	validPassword, err := server.roomService.IsValidPassword(context.Background(), roomID, password)
	if err != nil {
		wsSession.CloseWithMsg(melody.FormatCloseMessage(500, "Error: "+err.Error()))
		return
	}
	if !validPassword {
		wsSession.CloseWithMsg(melody.FormatCloseMessage(400, "Invalid password"))
		return
	}
	server.joinRoom(wsSession, roomID, userName)
}

func (server *HttpServer) sendAuthRequiredMessage(wsSession *melody.Session) error {
	return wsSession.Write([]byte("This room is protected, please enter the password"))
}

func (server *HttpServer) joinRoom(wsSession *melody.Session, roomID RoomID, userName string) {
	err := server.initializeChatSession(wsSession, roomID, userName)
	if err != nil {
		wsSession.CloseWithMsg(melody.FormatCloseMessage(500, "Error: "+err.Error()))
		return
	}

	if err := server.roomService.BroadcastConnectMessage(context.Background(), roomID, userName); err != nil {
		wsSession.CloseWithMsg(melody.FormatCloseMessage(500, "Error: "+err.Error()))
		return
	}
}

func (server *HttpServer) initializeChatSession(wsSession *melody.Session, roomID RoomID, userName string) error {
	ctx := context.Background()
	if err := server.roomService.AddRoomSubscriber(ctx, roomID, userName, server.msgSubscriber.topic); err != nil {
		return err
	}
	wsSession.Set(sessRidKey, roomID)
	return nil
}

func extractPassword(msg []byte) (string, error) {
	auth, err := decodeToRoomAuth(msg)
	if err != nil {
		return "", err
	}
	return auth.Password, nil
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
