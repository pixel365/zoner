package server

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pixel365/goepp/command"
	"github.com/pixel365/goepp/command/login"
	"github.com/pixel365/goepp/response"

	"github.com/pixel365/zoner/internal/repository/auth"

	conn2 "github.com/pixel365/zoner/epp/server/conn"
	"github.com/pixel365/zoner/internal/observability/metrics"
)

//nolint:gocognit,gocyclo,cyclop
func handleLogin(
	ctx context.Context,
	connection *conn2.Connection,
	cmd command.Commander,
	e *Epp,
) error {
	creds, ok := cmd.(*login.Login)
	if !ok {
		e.Metrics.Inc(ctx, metrics.AuthFailureTotal)
		e.Log.WithSessionId(connection.SessionId()).
			Error("cast failed", errors.New("invalid credentials type"))

		errorResponse := response.AnyError(response.CodeCommandUseError, response.CommandUseError)
		if err := connection.Write(ctx, errorResponse, e.Metrics.IncBytes); err != nil {
			return fmt.Errorf("write error response for invalid login command: %w", err)
		}
		return nil
	}

	creds.NewPassword = strings.TrimSpace(creds.NewPassword)

	if connection.IsAuthenticated() && creds.NewPassword == "" {
		errorResponse := response.AnyError(response.CodeObjectExists, "Already logged in")
		if err := connection.Write(ctx, errorResponse, e.Metrics.IncBytes); err != nil {
			return fmt.Errorf("write error response when client is authenticated: %w", err)
		}
		return nil
	}

	if creds.NewPassword != "" && len(creds.NewPassword) < e.Config.MinPasswordLength {
		errorResponse := response.AnyError(
			response.CodeParameterValueSyntaxError,
			fmt.Sprintf(
				"password length must be at least %d characters",
				e.Config.MinPasswordLength,
			),
		)
		if err := connection.Write(ctx, errorResponse, e.Metrics.IncBytes); err != nil {
			return fmt.Errorf("write error response for invalid password length: %w", err)
		}
		return nil
	}

	userId, maxActiveSessions, err := e.AuthService.Login(ctx, creds)
	if err != nil {
		var (
			errCode = 2200
			errType = response.AuthenticationError
		)

		switch {
		case errors.Is(err, auth.ErrInvalidCredentials),
			errors.Is(err, auth.ErrRegistrarIsBlocked),
			errors.Is(err, auth.ErrCannotChangePassword):
			e.Metrics.Inc(ctx, metrics.AuthFailureTotal)
		default:
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

	if creds.NewPassword == "" {
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
