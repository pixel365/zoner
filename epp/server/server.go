package server

import (
	"crypto/tls"

	"github.com/pixel365/zoner/epp/config"
	"github.com/pixel365/zoner/internal/logger"
	"github.com/pixel365/zoner/internal/repository"
	"github.com/pixel365/zoner/internal/repository/auth"
)

type Epp struct {
	Log            logger.Logger
	AuthRepository repository.Authenticator
	TLSConfig      *tls.Config
	Config         config.Epp
}

func NewEpp(cfg *config.Config, log logger.Logger) *Epp {
	cert, err := tls.LoadX509KeyPair(cfg.TLS.CertPath, cfg.TLS.KeyPath)
	if err != nil {
		return nil
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
	}
}
