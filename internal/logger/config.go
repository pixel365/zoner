package logger

import "os"

type ConfigOption func(*Config)

type Service struct {
	Name        string
	Version     string
	Environment string
	InstanceID  string
}

type Config struct {
	ECSVersion   string
	Service      Service
	EventDataset string
	LoggerName   string
}

func NewConfig(opts ...ConfigOption) *Config {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func WithECSVersion(version string) ConfigOption {
	return func(cfg *Config) {
		if version == "" {
			version = "8.11.0"
		}
		cfg.ECSVersion = version
	}
}

func WithService(service Service) ConfigOption {
	return func(cfg *Config) {
		dev := "dev"
		if service.Name == "" {
			serviceName := os.Getenv("SERVICE_NAME")
			if serviceName == "" {
				serviceName = "unknown"
			}
			service.Name = serviceName
		}

		if cfg.Service.Version == "" {
			version := os.Getenv("CI_BUILD_TAG")
			if version == "" {
				version = "0.0.1"
			}
			service.Version = version
		}

		if cfg.Service.Environment == "" {
			environment := os.Getenv("TARGET_SYSTEM")
			if environment == "" {
				environment = dev
			}
			service.Environment = environment
		}

		if cfg.Service.InstanceID == "" {
			instanceId := os.Getenv("POD_NAME")
			if instanceId == "" {
				instanceId = dev
			}
			service.InstanceID = instanceId
		}

		cfg.Service = service
	}
}

func WithEventDataset(dataset string) ConfigOption {
	return func(cfg *Config) {
		cfg.EventDataset = dataset
	}
}
