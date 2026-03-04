package server

import (
	"context"

	"github.com/pixel365/goepp/command"

	"github.com/pixel365/zoner/epp/server/conn"
)

func handleRenew(
	ctx context.Context,
	connection *conn.Connection,
	cmd command.Commander,
	e *Epp,
) error {
	return e.DomainService.Renew(ctx)
}
