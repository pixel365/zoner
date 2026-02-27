package otel

import (
	"context"
)

func (o *OtelMeter) Shutdown(ctx context.Context) error {
	var err error

	o.once.Do(func() {
		err = shutdownFn(ctx)
	})

	return err
}
