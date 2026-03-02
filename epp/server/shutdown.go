package server

import "context"

func (e *Epp) Shutdown(_ context.Context) {
	if e.DbPool != nil {
		e.DbPool.Close()
	}
}
