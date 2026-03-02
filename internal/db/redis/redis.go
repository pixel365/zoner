package redis

import (
	"context"

	"github.com/redis/go-redis/extra/redisotel/v9"
	r "github.com/redis/go-redis/v9"

	"github.com/pixel365/zoner/internal/repository"
)

var _ repository.SessionLimiter = (*Client)(nil)

type Client struct {
	client *r.Client
}

func (c *Client) Close() error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}

func NewRedisClient(ctx context.Context, cfg Config) *Client {
	if !cfg.IsValid() {
		panic("invalid redis config")
	}

	opts := &r.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
	}

	client := r.NewClient(opts)

	if err := redisotel.InstrumentTracing(client); err != nil {
		panic(err)
	}

	if err := client.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	c := &Client{client}

	return c
}
