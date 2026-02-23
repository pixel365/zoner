package server

import (
	"crypto/tls"
	"fmt"

	"github.com/pixel365/zoner/epp/config"
	"github.com/pixel365/zoner/internal/logger"
	"github.com/pixel365/zoner/internal/metrics"
	"github.com/pixel365/zoner/internal/repository"
	"github.com/pixel365/zoner/internal/repository/auth"
)

type Epp struct {
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

	return &Epp{
		Log:    log,
		Config: cfg.Epp,
		TLSConfig: &tls.Config{
			MinVersion:   tls.VersionTLS12,
			MaxVersion:   tls.VersionTLS13,
			Certificates: []tls.Certificate{cert},
		},
		AuthRepository: auth.NewAuth(),
		Metrics:        collector,
	}
}
