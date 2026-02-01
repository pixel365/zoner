package config

import (
	"crypto/tls"
	"errors"

	"github.com/pixel365/zoner/internal/logger"
)

type ServerOption func(*Config)

type Config struct {
	Log       logger.Logger
	TLSConfig *tls.Config
	Addr      string
}

func (c *Config) Validate() error {
	if c.Log == nil {
		return errors.New("logger is empty")
	}

	if c.TLSConfig == nil {
		return errors.New("tls config is empty")
	}

	if c.Addr == "" {
		return errors.New("address is empty")
	}

	return nil
}

func NewConfig(opts ...ServerOption) *Config {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func WithLogger(log logger.Logger) ServerOption {
	return func(cfg *Config) {
		cfg.Log = log
	}
}

func WithAddr(addr string) ServerOption {
	return func(cfg *Config) {
		cfg.Addr = addr
	}
}

func WithTLSConfig(tlsConfig *tls.Config) ServerOption {
	return func(cfg *Config) {
		cfg.TLSConfig = tlsConfig
	}
}
