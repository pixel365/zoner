package postgres

import (
	"context"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	log := cfg.Log
	config, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		if log != nil {
			log.Error("cannot parse config", err)
		}
		return nil, err
	}

	config.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		if log != nil {
			log.Error("cannot create pool", err)
		}
		return nil, err
	}

	if err = otelpgx.RecordStats(pool); err != nil {
		if log != nil {
			log.Error("cannot record stats", err)
		}
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		if log != nil {
			log.Error("cannot ping pool", err)
		}

		pool.Close()
		return nil, err
	}

	return pool, nil
}
