package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/pixel365/zoner/epp/config"
	"github.com/pixel365/zoner/epp/server"
	"github.com/pixel365/zoner/internal/logger"
	"github.com/pixel365/zoner/internal/logger/zerolog"
)

const appName = "epp-service"

var (
	tlsServerName string
	certPath      string
	certKeyPath   string
	eppListenAddr string
)

func init() {
	_ = godotenv.Load()

	tlsServerName = os.Getenv("TLS_SERVER_NAME")
	if tlsServerName == "" {
		tlsServerName = "localhost"
	}

	certPath = os.Getenv("CERT_PATH")
	certKeyPath = os.Getenv("CERT_KEY_PATH")

	eppListenAddr = os.Getenv("EPP_LISTEN_ADDR")
	if eppListenAddr == "" {
		eppListenAddr = ":7000"
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := zerolog.NewLogger(
		logger.NewConfig(
			logger.WithECSVersion("8.11.0"),
			logger.WithEventDataset(fmt.Sprintf("%s.logs", appName)),
			logger.WithService(logger.Service{
				Name: appName,
			}),
		),
	)

	mainLog := log.Component("main")

	if certPath == "" {
		mainLog.Error("cert error", errors.New("cert path is empty"))
		return
	}

	if certKeyPath == "" {
		mainLog.Error("cert error", errors.New("cert key path is empty"))
		return
	}

	cert, err := tls.LoadX509KeyPair(certPath, certKeyPath)
	if err != nil {
		mainLog.Error("load cert error", err)
		return
	}

	cfg := config.NewConfig(
		config.WithLogger(log),
		config.WithAddr(eppListenAddr),
		config.WithTLSConfig(&tls.Config{
			ServerName:   tlsServerName,
			MinVersion:   tls.VersionTLS12,
			MaxVersion:   tls.VersionTLS13,
			Certificates: []tls.Certificate{cert},
		}),
	)

	if err := cfg.Validate(); err != nil {
		mainLog.Error("config validation error", err)
		return
	}

	srv := server.NewEpp(cfg)
	if err := srv.Start(ctx); err != nil {
		mainLog.Error("epp server starting error", err)
	}
}
