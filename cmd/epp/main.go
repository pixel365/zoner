package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/pixel365/zoner/internal/metrics/noop"

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

	metrics := &noop.Noop{}

	srv := server.MustEpp(cfg, log, metrics)
	if err := srv.Start(ctx); err != nil {
		mainLog.Error("epp server starting error", err)
	}
}
