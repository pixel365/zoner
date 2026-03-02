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

	redisClient := redis.NewRedisClient(ctx, redis.NewConfigFromEnv())

	cfg.DB = pgPool
	cfg.RedisClient = redisClient

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
	go func() {
		if err := healthServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			mainLog.Error("health server start error", err)
		}
	}()

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := healthServer.Shutdown(shutdownCtx); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			mainLog.Error("health server shutdown error", err)
		}
	}()

	srv := server.MustEpp(cfg, log, metrics)
	defer srv.Shutdown(context.Background())

	if err := srv.Start(ctx, healthState.SetReady); err != nil {
		mainLog.Error("epp server starting error", err)
		stop()
	}
}
