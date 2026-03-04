package server

import (
	"context"
	"errors"

	"github.com/pixel365/goepp/command"

	"github.com/pixel365/zoner/epp/server/conn"
)

func handleUpdate(
	ctx context.Context,
	connection *conn.Connection,
	cmd command.Commander,
	e *Epp,
) error {
	return errors.New("not implemented")
}
