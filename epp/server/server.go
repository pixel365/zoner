package server

import (
	"crypto/tls"

	"github.com/pixel365/zoner/epp/config"
	"github.com/pixel365/zoner/internal/logger"
)

type Epp struct {
	Log       logger.Logger
	TLSConfig *tls.Config
	Addr      string
}

func NewEpp(config *config.Config) *Epp {
	return &Epp{
		Log:       config.Log,
		TLSConfig: config.TLSConfig,
		Addr:      config.Addr,
	}
}
