package middleware

import (
	"time"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ZapLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		dur := time.Since(start)
		status := c.Writer.Status()
		rid, _ := c.Get("request_id")
		log.Info("http_request",
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client", c.ClientIP()),
			zap.Duration("duration", dur),
			zap.Any("request_id", rid),
		)
	}
}