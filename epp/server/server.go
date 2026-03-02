package server

import (
	"crypto/tls"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pixel365/zoner/epp/config"
	"github.com/pixel365/zoner/internal/logger"
	"github.com/pixel365/zoner/internal/observability/metrics"
	"github.com/pixel365/zoner/internal/repository"
	"github.com/pixel365/zoner/internal/repository/auth"
)

type Epp struct {
	DbPool         *pgxpool.Pool
	Limiter        repository.SessionLimiter
	Log            logger.Logger
	AuthRepository repository.Authenticator
	Metrics        metrics.Collector
	TLSConfig      *tls.Config
	Config         config.Epp
}

func MustEpp(cfg *config.Config, log logger.Logger, collector metrics.Collector) *Epp {
	cert, err := tls.LoadX509KeyPair(cfg.TLS.CertPath, cfg.TLS.KeyPath)
	if err != nil {
		panic(fmt.Errorf("load cert and key: %w", err))
	}

	if collector == nil {
		panic("metrics collector is required")
	}

	if cfg.DB == nil {
		panic("db config is required")
	}

	if cfg.RedisClient == nil {
		panic("redis client is required")
	}

	return &Epp{
		DbPool:  cfg.DB,
		Limiter: cfg.RedisClient,
		Log:     log,
		Config:  cfg.Epp,
		TLSConfig: &tls.Config{
			MinVersion:   tls.VersionTLS12,
			MaxVersion:   tls.VersionTLS13,
			Certificates: []tls.Certificate{cert},
		},
		AuthRepository: auth.NewAuth(cfg.DB),
		Metrics:        collector,
	}
}
