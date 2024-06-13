package room

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/omran95/chatroom/pkg/common"
)

func NewRateLimiterMiddleware(createRoomsrateLimiter common.RateLimiter) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{createRoomsrateLimiter: createRoomsrateLimiter}
}

type RateLimiterMiddleware struct {
	createRoomsrateLimiter common.RateLimiter
}

func (rl *RateLimiterMiddleware) LimitCreateRooms(c *gin.Context) {
	hostIP := c.ClientIP()
	key := hostIP + ":create_room"

	allowed, err := rl.createRoomsrateLimiter.Allow(c.Request.Context(), key)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !allowed {
		c.Header("Retry-After", strconv.FormatFloat(1/rl.createRoomsrateLimiter.FillingRate, 'f', -1, 64))
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	c.Next()
}
