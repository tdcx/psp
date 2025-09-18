package middleware

import (
	"net"
	"sync"
	"time"
	"github.com/gin-gonic/gin"
)

type bucket struct {
	Tokens int
	Last   time.Time
}

func RateLimit(rps, burst int) gin.HandlerFunc {
	var mu sync.Mutex
	store := map[string]*bucket{}
	replenish := time.Duration(1e9 / max(1, rps))

	return func(c *gin.Context) {
		ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
		if ip == "" { ip = c.ClientIP() }
		mu.Lock()
		b, ok := store[ip]
		if !ok {
			b = &bucket{Tokens: burst, Last: time.Now()}
			store[ip] = b
		}
		elapsed := time.Since(b.Last)
		add := int(elapsed / replenish)
		if add > 0 {
			b.Tokens = min(burst, b.Tokens+add)
			b.Last = b.Last.Add(time.Duration(add) * replenish)
		}
		if b.Tokens <= 0 {
			mu.Unlock()
			c.AbortWithStatus(429)
			return
		}
		b.Tokens--
		mu.Unlock()
		c.Next()
	}
}

func min(a,b int) int { if a<b {return a}; return b }
func max(a,b int) int { if a>b {return a}; return b }