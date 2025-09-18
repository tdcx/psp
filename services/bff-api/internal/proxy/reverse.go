package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"github.com/gin-gonic/gin"
)

type Target struct { Base *url.URL; Prefix string }

func NewReverseProxy(raw, prefix string) (*httputil.ReverseProxy, *Target, error) {
	u, err := url.Parse(raw)
	if err != nil { return nil, nil, err }
	t := &Target{ Base: u, Prefix: prefix }
	rp := httputil.NewSingleHostReverseProxy(u)
	rp.ModifyResponse = func(resp *http.Response) error { return nil }
	return rp, t, nil
}

// Handler returns a Gin handler that strips the prefix and proxies to the target.
func Handler(rp *httputil.ReverseProxy, t *Target) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := strings.TrimPrefix(c.Request.URL.Path, t.Prefix)
		if p == c.Request.URL.Path { // prefix mismatch
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if !strings.HasPrefix(p, "/") { p = "/" + p }
		c.Request.URL.Path = p
		c.Request.Host = t.Base.Host
		rp.ServeHTTP(c.Writer, c.Request)
	}
}