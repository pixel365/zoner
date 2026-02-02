package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/pixel365/zoner/internal/logger"
)

type Config struct {
	LogLevel logger.LogLevel `yaml:"logLevel"`
	TLS      TLS             `yaml:"tls"`
	Epp      Epp             `yaml:"epp"`
}

func (c Config) Validate() error {
	if c.LogLevel == "" {
		return fmt.Errorf("log level is empty")
	}

	if err := c.TLS.Validate(); err != nil {
		return err
	}

	if err := c.Epp.Validate(); err != nil {
		return err
	}

	return nil
}

func MustConfig() *Config {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		panic("empty config path")
	}

	config, err := readConfig(path)
	if err != nil {
		panic(err)
	}

	return config
}

func readConfig(path string) (*Config, error) {
	var filePath string
	var err error

	filePath, err = filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	if !isValidPath(filePath, path) {
		return nil, fmt.Errorf("invalid config path: %s", filePath)
	}

	data, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func isValidPath(filePath, basePath string) bool {
	absBasePath, _ := filepath.Abs(basePath)
	absFilePath, _ := filepath.Abs(filePath)

	return strings.HasPrefix(absFilePath, absBasePath)
}
