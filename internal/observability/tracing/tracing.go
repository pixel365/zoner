package tracing

import "context"

type Span string
type Scope string

type KV struct {
	Key, Value string
}

type Tracer interface {
	Start(context.Context, Scope, Span, ...KV) (context.Context, Ender)
}

type Ender interface {
	End()
}
