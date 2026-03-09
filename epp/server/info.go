package server

import (
	"context"
	"fmt"

	"github.com/pixel365/goepp/command"
	"github.com/pixel365/goepp/command/info"
	"github.com/pixel365/goepp/response"

	"github.com/pixel365/zoner/epp/server/conn"
)

func handleInfo(
	ctx context.Context,
	connection *conn.Connection,
	cmd command.Commander,
	e *Epp,
) error {
	var resp response.Marshaller
	data, _ := cmd.(*info.Info)

	switch {
	case data.Domain != nil:
		return e.DomainService.Info(ctx)
	case data.Contact != nil:
		resp = e.ContactService.Info(ctx, *data.Contact, connection.UserId(), connection.Username())
	}

	if err := connection.Write(ctx, resp, e.Metrics.IncBytes); err != nil {
		return fmt.Errorf("write response error: %w", err)
	}

	return nil
}
