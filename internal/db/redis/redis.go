package redis

import (
	"context"

	"github.com/redis/go-redis/extra/redisotel/v9"
	r "github.com/redis/go-redis/v9"
)

func MustRedisClient(ctx context.Context, cfg Config) *r.Client {
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

	return client
}
