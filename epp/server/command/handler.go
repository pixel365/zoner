package command

import (
	"context"

	command2 "github.com/pixel365/zoner/epp/server/command/command"
)

type Handler interface {
	Handle(context.Context, command2.Commander)
}
