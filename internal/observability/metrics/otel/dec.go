package otel

import (
	"context"

	"github.com/pixel365/zoner/internal/observability/metrics"
)

func (o *OtelMeter) Dec(ctx context.Context, label metrics.Label, _ ...any) {
	if label == metrics.ActiveConnections {
		o.upDownCounter[label].Add(ctx, -1)
	}
}
