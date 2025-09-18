package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port             int           `envconfig:"PORT" default:"3000"`
	UpstreamSIF      string        `envconfig:"UPSTREAM_SIF" required:"true"`
	UpstreamPayments string        `envconfig:"UPSTREAM_PAYMENTS" required:"true"`
	CORSAllowOrigins string        `envconfig:"CORS_ALLOW_ORIGINS" default:"*"`
	RateLimitRPS     int           `envconfig:"RATE_LIMIT_RPS" default:"10"`
	RateLimitBurst   int           `envconfig:"RATE_LIMIT_BURST" default:"20"`
	LogLevel         string        `envconfig:"LOG_LEVEL" default:"info"`
	ShutdownTimeout  time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"10s"`
}

func Load() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	return cfg, err
}
