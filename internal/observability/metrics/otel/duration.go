package otel

import (
	"context"
	"time"

	"github.com/pixel365/zoner/internal/observability/metrics"
)

func (o *OtelMeter) Duration(
	ctx context.Context,
	label metrics.Label,
	duration time.Duration,
	_ ...any,
) {
	if label == metrics.SessionDurationMs {
		o.f64Histogram[label].Record(ctx, float64(duration.Milliseconds()))
	}
}
