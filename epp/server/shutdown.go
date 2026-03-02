package server

import (
	"context"
	"io"
)

func (e *Epp) Shutdown(_ context.Context) {
	if e.DbPool != nil {
		e.DbPool.Close()
	}

	if e.Limiter != nil {
		if c, ok := e.Limiter.(io.Closer); ok {
			_ = c.Close()
		}
	}
}
