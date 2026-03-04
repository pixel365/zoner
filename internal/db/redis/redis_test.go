package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRedisClient(t *testing.T) {
	assert.Panics(t, func() {
		MustRedisClient(context.Background(), Config{})
	})

	assert.Panics(t, func() {
		MustRedisClient(context.Background(), Config{
			Host: "localhost",
			Port: "6379",
		})
	})
}
