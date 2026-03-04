package server

import (
	"crypto/tls"
	"fmt"

	auth2 "github.com/pixel365/zoner/internal/service/auth"
	"github.com/pixel365/zoner/internal/service/contact"
	"github.com/pixel365/zoner/internal/service/limiter"
	"github.com/pixel365/zoner/internal/service/zone"

	"github.com/pixel365/zoner/internal/service/domain"

	"github.com/pixel365/zoner/epp/config"
	"github.com/pixel365/zoner/internal/logger"
	"github.com/pixel365/zoner/internal/observability/metrics"
)

type Epp struct {
	Log            logger.Logger
	Metrics        metrics.Collector
	LimiterService limiter.LimiterService
	AuthService    auth2.AuthService
	DomainService  domain.DomainService
	ContactService contact.ContactService
	ZoneService    zone.ZoneService
	TLSConfig      *tls.Config
	Config         config.Epp
}

func MustEpp(
	cfg *config.Config,
	log logger.Logger,
	collector metrics.Collector,
	limiterSvc limiter.LimiterService,
	authSvc auth2.AuthService,
	domainSvc domain.DomainService,
	contactSvc contact.ContactService,
	zoneSvc zone.ZoneService,
) *Epp {
	cert, err := tls.LoadX509KeyPair(cfg.TLS.CertPath, cfg.TLS.KeyPath)
	if err != nil {
		panic(fmt.Errorf("load cert and key: %w", err))
	}

	if collector == nil {
		panic("metrics collector is required")
	}

	if limiterSvc == nil {
		panic("limiter service is required")
	}

	if authSvc == nil {
		panic("auth service is required")
	}

	if domainSvc == nil {
		panic("domain service is required")
	}

	if contactSvc == nil {
		panic("contact service is required")
	}

	if zoneSvc == nil {
		panic("zone service is required")
	}

	return &Epp{
		Log:    log,
		Config: cfg.Epp,
		TLSConfig: &tls.Config{
			MinVersion:   tls.VersionTLS12,
			MaxVersion:   tls.VersionTLS13,
			Certificates: []tls.Certificate{cert},
		},
		Metrics:        collector,
		LimiterService: limiterSvc,
		AuthService:    authSvc,
		DomainService:  domainSvc,
		ContactService: contactSvc,
		ZoneService:    zoneSvc,
	}
}
