package otel

import (
	"context"

	"github.com/pixel365/zoner/internal/observability/metrics"
)

func (o *OtelMeter) IncBytes(
	ctx context.Context,
	label metrics.Label,
	size int64,
	_ ...any,
) {
	//nolint:exhaustive
	switch label {
	case metrics.FramesReadTotal, metrics.FramesWriteTotal:
		o.int64Histogram[label].Record(ctx, size)
	default:
		return
	}
}
