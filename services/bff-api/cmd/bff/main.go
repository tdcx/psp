package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tdcx/psp/services/bff-api/internal/config"
	"github.com/tdcx/psp/services/bff-api/internal/server"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	logCfg := zap.NewProductionConfig()
	if cfg.LogLevel == "debug" {
		logCfg = zap.NewDevelopmentConfig()
	}
	logger, _ := logCfg.Build()
	defer logger.Sync()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	logger.Info("starting bff", zap.Int("port", cfg.Port))
	if err := server.New(cfg, logger).Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
	}
}
