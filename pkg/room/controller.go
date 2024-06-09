package room

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omran95/chat-app/pkg/common"
)

func (server *HttpServer) CreateRoom(c *gin.Context) {
	room, err := server.roomService.CreateRoom(c)
	if err != nil {
		server.logger.Error(err.Error())
		response(c, http.StatusInternalServerError, common.ErrServer)
		return
	}
	c.JSON(http.StatusCreated, room)
}

func response(c *gin.Context, httpCode int, err error) {
	message := err.Error()
	c.JSON(httpCode, common.ErrResponse{
		Message: message,
	})
}
