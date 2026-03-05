package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"

	"github.com/pixel365/zoner/internal/repository/limiter"
	limiter2 "github.com/pixel365/zoner/internal/service/limiter"

	"github.com/pixel365/zoner/internal/repository/auth"
	auth2 "github.com/pixel365/zoner/internal/service/auth"

	"github.com/pixel365/zoner/internal/repository/contact"
	"github.com/pixel365/zoner/internal/repository/domain"
	"github.com/pixel365/zoner/internal/repository/zone"
	contact2 "github.com/pixel365/zoner/internal/service/contact"
	domain2 "github.com/pixel365/zoner/internal/service/domain"
	zone2 "github.com/pixel365/zoner/internal/service/zone"

	"github.com/pixel365/zoner/internal/db/redis"

	"github.com/pixel365/zoner/internal/db/postgres"

	"github.com/pixel365/zoner/internal/health"
	"github.com/pixel365/zoner/internal/observability/metrics/collector"

	"github.com/pixel365/zoner/epp/config"
	"github.com/pixel365/zoner/epp/server"
	"github.com/pixel365/zoner/internal/logger"
	"github.com/pixel365/zoner/internal/logger/zerolog"
)

const appName = "epp-service"

func init() {
	_ = godotenv.Load()
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.MustConfig()

	log := zerolog.NewLogger(
		logger.NewConfig(
			logger.WithECSVersion("8.11.0"),
			logger.WithEventDataset(fmt.Sprintf("%s.logs", appName)),
			logger.WithService(logger.Service{
				Name: appName,
			}),
			logger.WithLogLevel(cfg.LogLevel),
		),
	)

	mainLog := log.Component("main")

	if err := cfg.Validate(); err != nil {
		mainLog.Error("config validation error", err)
		return
	}

	pgConfig := postgres.NewConfigFromEnv()
	pgConfig.Log = log.Component("postgres.db")
	pgPool, err := postgres.NewPool(ctx, pgConfig)
	if err != nil {
		mainLog.Error("postgres pool error", err)
		return
	}
	defer pgPool.Close()

	redisClient := redis.MustRedisClient(ctx, redis.NewConfigFromEnv())
	defer func() {
		_ = redisClient.Close()
	}()

	metrics := collector.NewCollector(ctx, log)
	defer func() {
		if err := metrics.Shutdown(context.Background()); err != nil {
			mainLog.Error("metrics shutdown error", err)
		}
	}()

	_ = runtime.Start(runtime.WithMinimumReadMemStatsInterval(10 * time.Second))
	_ = host.Start()

	healthState := health.NewState()
	healthServer := health.NewHealthServer(healthState)
	if healthServer != nil {
		go func() {
			if err := healthServer.ListenAndServe(); err != nil &&
				!errors.Is(err, http.ErrServerClosed) {
				mainLog.Error("health server start error", err)
			}
		}()
	}

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if healthServer != nil {
			if err := healthServer.Shutdown(shutdownCtx); err != nil &&
				!errors.Is(err, http.ErrServerClosed) {
				mainLog.Error("health server shutdown error", err)
			}
		}
	}()

	limiterRepo := limiter.NewRepository(redisClient)
	authRepo := auth.NewRepository(pgPool)
	domainRepo := domain.NewRepository(pgPool)
	contactRepo := contact.NewRepository(pgPool)
	zoneRepo := zone.NewRepository(pgPool)

	limiterSvc := limiter2.MustService(limiterRepo)
	authSvc := auth2.MustService(authRepo)
	domainSvc := domain2.MustService(domainRepo)
	contactSvc := contact2.MustService(contactRepo, log)
	zoneSvc := zone2.MustService(zoneRepo)

	srv := server.MustEpp(cfg, log, metrics, limiterSvc, authSvc, domainSvc, contactSvc, zoneSvc)
	defer srv.Shutdown(context.Background())

	if err := srv.Start(ctx, healthState.SetReady); err != nil {
		mainLog.Error("epp server starting error", err)
		stop()
	}
}
