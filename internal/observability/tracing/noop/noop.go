package noop

import (
	"context"

	"github.com/pixel365/zoner/internal/observability/tracing"
)

var _ tracing.Tracer = (*Noop)(nil)
var _ tracing.Ender = (*Noop)(nil)

type Noop struct{}

func (n *Noop) Start(
	context.Context,
	tracing.Scope,
	tracing.Span,
	...tracing.KV,
) (context.Context, tracing.Ender) {
	return context.TODO(), n
}

func (n *Noop) End() {}
