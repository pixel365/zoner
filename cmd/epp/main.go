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

	metrics := collector.NewCollector(ctx, log)
	defer func() {
		if err := metrics.Shutdown(context.Background()); err != nil {
			mainLog.Error("metrics shutdown error", err)
		}
	}()

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
	if err := srv.Start(ctx, healthState.SetReady); err != nil {
		mainLog.Error("epp server starting error", err)
		stop()
	}
}
