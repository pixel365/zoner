package config

import "errors"

type TLS struct {
	ServerName string `yaml:"serverName"`
	CertPath   string `yaml:"certPath"`
	KeyPath    string `yaml:"keyPath"`
}

func (t TLS) Validate() error {
	if t.ServerName == "" {
		return errors.New("server name is empty")
	}

	if t.CertPath == "" {
		return errors.New("cert path is empty")
	}

	if t.KeyPath == "" {
		return errors.New("key path is empty")
	}

	return nil
}
