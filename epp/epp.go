package epp

import (
	"context"
)

type Server interface {
	Start(context.Context) error
	Shutdown(context.Context) error
}

type Conn interface {
	ReadFrame(context.Context) ([]byte, error)
	WriteFrame(context.Context, []byte) error
	Close() error
}
