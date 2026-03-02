package redis

import "os"

type ConfigOption func(*Config)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
}

func (c Config) IsValid() bool {
	return len(c.Password) > 0 && len(c.Host) > 0 && len(c.Port) > 0
}

func NewConfig(opts ...ConfigOption) Config {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return *cfg
}

func NewConfigFromEnv() Config {
	return NewConfig(
		WithUsername(os.Getenv("REDIS_USERNAME")),
		WithPassword(os.Getenv("REDIS_PASSWORD")),
		WithHost(os.Getenv("REDIS_HOST")),
		WithPort(os.Getenv("REDIS_PORT")),
	)
}

func WithUsername(username string) ConfigOption {
	return func(cfg *Config) {
		cfg.Username = username
	}
}

func WithPassword(password string) ConfigOption {
	return func(cfg *Config) {
		cfg.Password = password
	}
}

func WithHost(host string) ConfigOption {
	return func(cfg *Config) {
		cfg.Host = host
	}
}

func WithPort(port string) ConfigOption {
	return func(cfg *Config) {
		cfg.Port = port
	}
}
