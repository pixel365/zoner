package server

import (
	"context"
	"fmt"

	"github.com/pixel365/goepp/command/greeting"

	conn2 "github.com/pixel365/zoner/epp/server/conn"
)

func handleHello(ctx context.Context, connection *conn2.Connection, e *Epp) error {
	g := greeting.NewGreeting(e.Config.Greeting)
	if err := connection.Write(ctx, g, e.Metrics.IncBytes); err != nil {
		return fmt.Errorf("write greeting error: %w", err)
	}
	return nil
}
