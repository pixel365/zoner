package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig(
		WithHost("localhost"),
		WithPort("6379"),
		WithPassword("password"),
		WithUsername("username"),
	)

	assert.NotNil(t, cfg)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "6379", cfg.Port)
	assert.Equal(t, "password", cfg.Password)
	assert.Equal(t, "username", cfg.Username)
}
