package otel

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"

	"github.com/pixel365/zoner/internal/observability/metrics"
)

var _ metrics.Collector = (*OtelMeter)(nil)

const serviceName = "zoner"

var errInitialization error

type OtelMeter struct {
	int64Counter   map[metrics.Label]metric.Int64Counter
	f64Histogram   map[metrics.Label]metric.Float64Histogram
	int64Histogram map[metrics.Label]metric.Int64Histogram
	upDownCounter  map[metrics.Label]metric.Int64UpDownCounter
	once           sync.Once
}

func NewMeter(ctx context.Context) (*OtelMeter, error) {
	once.Do(func() {
		if err := initMeter(ctx, serviceName); err != nil {
			errInitialization = err
			return
		}
	})

	if errInitialization != nil {
		return nil, errInitialization
	}

	o := OtelMeter{
		int64Counter:   make(map[metrics.Label]metric.Int64Counter),
		f64Histogram:   make(map[metrics.Label]metric.Float64Histogram),
		int64Histogram: make(map[metrics.Label]metric.Int64Histogram),
		upDownCounter:  make(map[metrics.Label]metric.Int64UpDownCounter),
	}

	var int64Requests metric.Int64Counter
	var f64HistogramRequests metric.Float64Histogram
	var int64HistogramRequests metric.Int64Histogram
	var upDownCounterRequests metric.Int64UpDownCounter

	meter := otel.Meter(serviceName)

	int64Counters := map[metrics.Label]struct{}{
		metrics.ConnectionsTotal:        {},
		metrics.RequestsTotal:           {},
		metrics.CommandsTotal:           {},
		metrics.CommandsWithErrorsTotal: {},
		metrics.AuthSuccessTotal:        {},
		metrics.AuthFailureTotal:        {},
		metrics.ParseErrorsTotal:        {},
	}

	for label := range int64Counters {
		int64Requests, _ = meter.Int64Counter(
			fmt.Sprintf("%s.%s", serviceName, label),
			metric.WithDescription(label.Description()),
		)
		o.int64Counter[label] = int64Requests
	}

	f64Histograms := map[metrics.Label]struct{}{
		metrics.SessionDurationMs: {},
	}

	for label := range f64Histograms {
		f64HistogramRequests, _ = meter.Float64Histogram(
			fmt.Sprintf("%s.%s", serviceName, label),
			metric.WithDescription(label.Description()),
			metric.WithUnit("ms"),
			metric.WithExplicitBucketBoundaries(
				1, 2, 5, 10, 25, 50, 75, 100,
				250, 500, 750, 1000,
				1500, 2500, 5000, 7500, 10000,
			),
		)
		o.f64Histogram[label] = f64HistogramRequests
	}

	int64Histograms := map[metrics.Label]struct{}{
		metrics.FramesReadTotal:  {},
		metrics.FramesWriteTotal: {},
	}

	for label := range int64Histograms {
		int64HistogramRequests, _ = meter.Int64Histogram(
			fmt.Sprintf("%s.%s", serviceName, label),
			metric.WithDescription(label.Description()),
			metric.WithUnit("B"),
			metric.WithExplicitBucketBoundaries(
				0,
				128,
				512,
				1<<10,
				4<<10,
				16<<10,
				64<<10,
				256<<10,
				1<<20,
				4<<20,
			),
		)
		o.int64Histogram[label] = int64HistogramRequests
	}

	upDownCounters := map[metrics.Label]struct{}{
		metrics.ActiveConnections: {},
	}

	for label := range upDownCounters {
		upDownCounterRequests, _ = meter.Int64UpDownCounter(
			fmt.Sprintf("%s.%s", serviceName, label),
			metric.WithDescription(label.Description()),
		)
		o.upDownCounter[label] = upDownCounterRequests
	}

	return &o, nil
}
