package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// AddCurrentTime adds the current time to the context
func AddCurrentTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("currentTime", time.Now())
		c.Next()
	}
}
