package epp

import (
	"context"
)

type Server interface {
	Start(context.Context, func(bool)) error
	Shutdown(context.Context)
}

type Conn interface {
	ReadFrame(context.Context) ([]byte, error)
	WriteFrame(context.Context, []byte) error
	Close() error
}
