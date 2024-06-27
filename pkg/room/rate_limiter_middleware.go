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
	tokensPerRequest := 10
	allowed, retryAfter, err := rl.createRoomsrateLimiter.Allow(c.Request.Context(), key, tokensPerRequest)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !allowed {
		c.Header("Retry-After", strconv.FormatInt(int64(retryAfter), 10))
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	c.Next()
}
