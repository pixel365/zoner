package collector

import (
	"context"
	"os"
	"strconv"
	"sync"

	"github.com/pixel365/zoner/internal/logger"
	"github.com/pixel365/zoner/internal/observability/metrics"
	"github.com/pixel365/zoner/internal/observability/metrics/noop"
	"github.com/pixel365/zoner/internal/observability/metrics/otel"
)

var (
	enabled bool
	once    sync.Once
)

func NewCollector(ctx context.Context, log logger.Logger) metrics.Collector {
	once.Do(func() {
		enabled, _ = strconv.ParseBool(os.Getenv("METRICS_ENABLED"))
	})

	if !enabled {
		return &noop.Noop{}
	}

	meter, err := otel.NewMeter(ctx)
	if err != nil {
		log.Error("otel initialization error", err)
		return &noop.Noop{}
	}

	return meter
}
