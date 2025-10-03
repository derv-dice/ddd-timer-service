package http

import (
	"ddd-timer-service/internal/pkg/tracelog"

	"github.com/gin-gonic/gin"
)

func accessLogMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		tl, ctx := tracelog.Begin(c.Request.Context(), "HTTP")
		defer tl.End()

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		tl.InfoWithDuration("request processed",
			tracelog.Int("status", c.Writer.Status()),
			tracelog.String("ip", c.ClientIP()),
			tracelog.String("method", c.Request.Method),
			tracelog.String("path", c.Request.URL.Path))
	}
}
