package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pixel365/goepp/command"
	"github.com/pixel365/goepp/command/login"
	"github.com/pixel365/goepp/response"

	conn2 "github.com/pixel365/zoner/epp/server/conn"
	errors2 "github.com/pixel365/zoner/internal/errors"
	"github.com/pixel365/zoner/internal/observability/metrics"
)

func handleLogin(
	ctx context.Context,
	connection *conn2.Connection,
	cmd command.Commander,
	e *Epp,
) error {
	if connection.IsAuthenticated() {
		errorResponse := response.AnyError(2302, "Already logged in")
		if err := connection.Write(ctx, errorResponse, e.Metrics.IncBytes); err != nil {
			return fmt.Errorf("write error response when client is authenticated: %w", err)
		}
		return nil
	}

	creds, ok := cmd.(*login.Login)
	if !ok {
		e.Metrics.Inc(ctx, metrics.AuthFailureTotal)
		e.Log.WithSessionId(connection.SessionId()).
			Error("cast failed", errors.New("invalid credentials type"))

		errorResponse := response.AnyError(2002, response.CommandUseError)
		if err := connection.Write(ctx, errorResponse, e.Metrics.IncBytes); err != nil {
			return fmt.Errorf("write error response for invalid login command: %w", err)
		}
		return nil
	}

	userId, maxActiveSessions, err := e.AuthService.Login(ctx, creds.ClientID, creds.Password)
	if err != nil {
		var (
			errCode = 2200
			errType = response.AuthenticationError
		)
		if errors.Is(err, errors2.ErrInvalidCredentials) {
			e.Metrics.Inc(ctx, metrics.AuthFailureTotal)
		} else {
			errCode = 2400
			errType = response.CommandFailed
		}

		e.Log.WithSessionId(connection.SessionId()).
			WithUsername(creds.ClientID).
			WithUserId(userId).
			Error("login failed", err)

		errorResponse := response.AnyError(response.ResponseCode(errCode), errType)
		if err = connection.Write(ctx, errorResponse, e.Metrics.IncBytes); err != nil {
			return fmt.Errorf("write error response for invalid login credentials: %w", err)
		}

		return nil
	}

	connection.SetClientUsername(creds.ClientID)
	connection.SetUserId(userId)

	reserved, err := e.LimiterService.Reserve(
		ctx,
		connection.SessionKey(),
		maxActiveSessions,
		time.Duration(e.Config.ActiveSessionTtl)*time.Second,
	)
	if err != nil {
		e.Log.WithSessionId(connection.SessionId()).
			WithUsername(creds.ClientID).
			WithUserId(userId).
			Error("session reserve failed", err)
		errorResponse := response.AnyError(2400, response.CommandFailed)
		if err = connection.Write(ctx, errorResponse, e.Metrics.IncBytes); err != nil {
			return fmt.Errorf("write error response for session reserve error: %w", err)
		}
		return nil
	}

	if !reserved {
		e.Log.WithSessionId(connection.SessionId()).
			WithUsername(creds.ClientID).
			WithUserId(userId).
			Info("session limit exceeded")
		errorResponse := response.AnyError(
			2502,
			response.SessionLimitExceededServerClosingConnection,
		)
		if err = connection.Write(ctx, errorResponse, e.Metrics.IncBytes); err != nil {
			return fmt.Errorf("write error response for session limit exceeded: %w", err)
		}
		return errSessionTerminate
	}

	connection.SetAuthenticated(true)
	res := response.NewResponse[struct{}, struct{}](1000, response.CommandCompletedSuccessfully)
	if err := connection.Write(ctx, res, e.Metrics.IncBytes); err != nil {
		return fmt.Errorf("write login response error: %w", err)
	}

	e.Metrics.Inc(ctx, metrics.AuthSuccessTotal)
	e.Log.WithSessionId(connection.SessionId()).
		WithUsername(creds.ClientID).
		WithUserId(userId).
		Info("login successful")

	return nil
}
