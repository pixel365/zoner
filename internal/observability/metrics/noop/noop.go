package noop

import (
	"context"
	"time"

	"github.com/pixel365/zoner/internal/observability/metrics"
)

var _ metrics.Collector = (*Noop)(nil)

type Noop struct{}

func (n Noop) Inc(context.Context, metrics.Label, ...any)                     {}
func (n Noop) Dec(context.Context, metrics.Label, ...any)                     {}
func (n Noop) Duration(context.Context, metrics.Label, time.Duration, ...any) {}
func (n Noop) IncBytes(context.Context, metrics.Label, int64, ...any)         {}

func (n Noop) Shutdown(
	context.Context,
) error {
	return nil
}
