package config

import (
	"errors"

	"github.com/pixel365/zoner/epp/config/epp/greeting"
)

type Epp struct {
	ListenAddr   string            `yaml:"listenAddr"`
	Greeting     greeting.Greeting `yaml:"greeting"`
	IdleTimeout  int               `yaml:"idleTimeout"`
	ReadTimeout  int               `yaml:"readTimeout"`
	WriteTimeout int               `yaml:"writeTimeout"`
}

func (e *Epp) Validate() error {
	if e.ListenAddr == "" {
		return errors.New("listen addr is empty")
	}

	if e.IdleTimeout <= 0 {
		return errors.New("idle timeout must be positive")
	}

	if e.ReadTimeout <= 0 {
		return errors.New("read timeout must be positive")
	}

	if e.WriteTimeout <= 0 {
		return errors.New("write timeout must be positive")
	}

	if err := e.Greeting.Validate(); err != nil {
		return err
	}

	return nil
}
