package server

import (
	"context"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/yourorg/psp/services/bff-api/internal/config"
	"github.com/yourorg/psp/services/bff-api/internal/middleware"
	"github.com/yourorg/psp/services/bff-api/internal/proxy"
)

type Server struct {
	cfg config.Config
	log *zap.Logger
}

func New(cfg config.Config, log *zap.Logger) *Server { return &Server{cfg: cfg, log: log} }

func (s *Server) Run(ctx context.Context) error {
	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(middleware.RequestID())
	g.Use(middleware.ZapLogger(s.log))
	g.Use(middleware.CORS(s.cfg.CORSAllowOrigins))
	g.Use(middleware.RateLimit(s.cfg.RateLimitRPS, s.cfg.RateLimitBurst))

	g.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// Reverse proxies (strip prefix)
	rpSIF, tgtSIF, err := proxy.NewReverseProxy(s.cfg.UpstreamSIF, "/api/sif")
	if err != nil { return err }
	g.Any("/api/sif/*any", proxy.Handler(rpSIF, tgtSIF))

	rpPay, tgtPay, err := proxy.NewReverseProxy(s.cfg.UpstreamPayments, "/api/payments")
	if err != nil { return err }
	g.Any("/api/payments/*any", proxy.Handler(rpPay, tgtPay))

	srv := &http.Server{ Addr: fmt.Sprintf(":%d", s.cfg.Port), Handler: g }
	go func(){ _ = srv.ListenAndServe() }()
	<-ctx.Done()
	ctxShut, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()
	return srv.Shutdown(ctxShut)
}