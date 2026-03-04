package server

import (
	"context"
	"fmt"

	conn2 "github.com/pixel365/zoner/epp/server/conn"
	"github.com/pixel365/zoner/epp/server/response"
)

func handleLogout(ctx context.Context, connection *conn2.Connection, e *Epp) error {
	if connection.IsAuthenticated() {
		res := response.NewResponse[struct{}, struct{}](
			1500,
			response.CommandCompleteSuccessfullyEndingSession,
		)
		if err := connection.Write(ctx, res, e.Metrics.IncBytes); err != nil {
			return fmt.Errorf("write logout response error: %w", err)
		}

		if err := e.LimiterService.Release(ctx, connection.SessionKey()); err != nil {
			e.Log.WithSessionId(connection.SessionId()).
				WithUserId(connection.UserId()).
				Error("session release failed", err)
		}

		connection.SetAuthenticated(false)

		return errSessionTerminate
	}

	errorResponse := response.AnyError(2002, response.CommandUseError)
	if err := connection.Write(ctx, errorResponse, e.Metrics.IncBytes); err != nil {
		return fmt.Errorf("write error response for invalid logout command: %w", err)
	}

	return nil
}
