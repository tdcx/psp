package middleware

import (
	"strings"
	"github.com/gin-gonic/gin"
)

func CORS(allowOrigins string) gin.HandlerFunc {
	allowed := strings.Split(allowOrigins, ",")
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		ok := allowOrigins == "*"
		if !ok {
			for _, a := range allowed {
				if strings.TrimSpace(a) == origin {
					ok = true
					break
				}
			}
		}
		if ok && origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Vary", "Origin")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-Id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}