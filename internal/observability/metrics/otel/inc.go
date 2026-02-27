package otel

import (
	"context"

	"github.com/pixel365/zoner/internal/observability/metrics"
)

func (o *OtelMeter) Inc(ctx context.Context, label metrics.Label, _ ...any) {
	//nolint:exhaustive
	switch label {
	case metrics.ConnectionsTotal,
		metrics.RequestsTotal,
		metrics.CommandsTotal,
		metrics.CommandsWithErrorsTotal,
		metrics.AuthSuccessTotal,
		metrics.AuthFailureTotal,
		metrics.ParseErrorsTotal:
		o.int64Counter[label].Add(ctx, 1)
	case metrics.ActiveConnections:
		o.upDownCounter[label].Add(ctx, 1)
	default:
		return
	}
}
