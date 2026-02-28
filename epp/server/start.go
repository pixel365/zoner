package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/pixel365/goepp"
	"github.com/pixel365/goepp/command"
	"github.com/pixel365/goepp/command/greeting"
	"github.com/pixel365/goepp/command/login"

	"github.com/pixel365/zoner/internal/observability/metrics"

	conn2 "github.com/pixel365/zoner/epp/server/conn"
	"github.com/pixel365/zoner/epp/server/response"
)

var errSessionTerminate = errors.New("session terminate")

func (e *Epp) Start(ctx context.Context, readyFn func(bool)) error {
	if readyFn == nil {
		return fmt.Errorf("ready func is nil")
	}

	ln, err := tls.Listen("tcp", e.Config.ListenAddr, e.TLSConfig)
	if err != nil {
		readyFn(false)
		return err
	}

	go func() {
		<-ctx.Done()
		readyFn(false)
		if err := ln.Close(); err != nil {
			e.Log.Error("close listener error", err)
		}
	}()

	readyFn(true)

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				continue
			}
		}

		e.Metrics.Inc(ctx, metrics.RequestsTotal)

		tlsConn, ok := conn.(*tls.Conn)
		if ok {
			if err := tlsConn.HandshakeContext(ctx); err != nil {
				if err = conn.Close(); err != nil {
					e.Log.Error("close connection error", err)
				}
				continue
			}
		}

		go e.handleConnection(ctx, conn)
	}
}

// TODO: rm gocognit
//
//nolint:gocognit
func (e *Epp) handleConnection(ctx context.Context, conn net.Conn) {
	e.Metrics.Inc(ctx, metrics.ConnectionsTotal)
	e.Metrics.Inc(ctx, metrics.ActiveConnections)

	connection := conn2.NewConnection(conn, &e.Config)
	log := e.Log.WithSessionId(connection.SessionId()).WithAddress(conn.RemoteAddr().String())

	log.Info("session.open")

	var closeError error

	defer func() {
		duration := time.Since(connection.SessionStart())

		e.Metrics.Dec(ctx, metrics.ActiveConnections)
		e.Metrics.Duration(ctx, metrics.SessionDurationMs, duration)

		if err := connection.Close(); err != nil {
			log.WithUserId(connection.UserId()).
				WithEventDuration(duration).
				Error("close connection error", err)
			return
		}

		e := &closeError
		if e != nil && *e != nil {
			log.WithUserId(connection.UserId()).
				WithEventDuration(duration).
				Error("session.close", *e)
			return
		}

		log.WithUserId(connection.UserId()).WithEventDuration(duration).Info("session.close")
	}()

	g := greeting.NewGreeting(e.Config.Greeting)
	if err := connection.Write(ctx, g, e.Metrics.IncBytes); err != nil {
		closeError = err
		log.Error("write greeting error", err)
		return
	}

	parser := goepp.CmdParser{}

	for {
		select {
		case <-ctx.Done():
			closeError = fmt.Errorf("context done")
			return
		default:
			frame, err := connection.ReadFrame(ctx, e.Metrics.IncBytes)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				closeError = fmt.Errorf("read frame error: %w", err)
				return
			}

			cmd, err := parseFrame(ctx, connection, &parser, frame, e)
			if err != nil {
				closeError = fmt.Errorf("parse frame error: %w", err)
				return
			}

			e.Metrics.Inc(ctx, metrics.CommandsTotal)

			err = sendResponse(ctx, connection, cmd, e)
			if errors.Is(err, errSessionTerminate) {
				return
			}

			if err != nil {
				closeError = fmt.Errorf("send response error: %w", err)
				e.Metrics.Inc(ctx, metrics.CommandsWithErrorsTotal)
				return
			}
		}
	}
}

func parseFrame(
	ctx context.Context,
	connection *conn2.Connection,
	parser *goepp.CmdParser,
	frame []byte,
	e *Epp,
) (command.Commander, error) {
	cmd, err := parser.Parse(frame)
	if err != nil {
		e.Metrics.Inc(ctx, metrics.ParseErrorsTotal)
		errorResponse := response.AnyError(2001, response.CommandSyntaxError)
		if err = connection.Write(ctx, errorResponse, e.Metrics.IncBytes); err != nil {
			return nil, fmt.Errorf("write error response for invalid command: %w", err)
		}
		return nil, nil
	}

	return cmd, nil
}

func sendResponse(
	ctx context.Context,
	connection *conn2.Connection,
	cmd command.Commander,
	e *Epp,
) error {
	if cmd == nil {
		return nil
	}

	if cmd.Name() == command.Hello {
		g := greeting.NewGreeting(e.Config.Greeting)
		if err := connection.Write(ctx, g, e.Metrics.IncBytes); err != nil {
			return fmt.Errorf("write greeting error: %w", err)
		}
		return nil
	}

	if cmd.Name() == command.Logout {
		return handleLogout(ctx, connection, e)
	}

	if cmd.Name() == command.Login {
		return handleLogin(ctx, connection, cmd, e)
	}

	if cmd.NeedAuth() && !connection.IsAuthenticated() {
		errorResponse := response.AnyError(2200, response.AuthorizationError)
		if err := connection.Write(ctx, errorResponse, e.Metrics.IncBytes); err != nil {
			return fmt.Errorf("write error response when client is not authenticated: %w", err)
		}
		return nil
	}

	errorResponse := response.AnyError(2101, response.UnimplementedCommand)
	if err := connection.Write(ctx, errorResponse, e.Metrics.IncBytes); err != nil {
		return fmt.Errorf("write error response for unimplemented command: %w", err)
	}

	return nil
}

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

	creds, ok := cmd.(login.Login)
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

	if err := e.AuthRepository.Login(creds.ClientID, creds.Password); err != nil {
		e.Metrics.Inc(ctx, metrics.AuthFailureTotal)
		e.Log.WithSessionId(connection.SessionId()).
			WithUserId(creds.ClientID).
			Error("login failed", err)

		errorResponse := response.AnyError(2201, response.AuthorizationError)
		if err = connection.Write(ctx, errorResponse, e.Metrics.IncBytes); err != nil {
			return fmt.Errorf("write error response for invalid login credentials: %w", err)
		}

		return nil
	}

	res := response.NewResponse[struct{}, struct{}](1000, response.CommandCompletedSuccessfully)
	if err := connection.Write(ctx, res, e.Metrics.IncBytes); err != nil {
		return fmt.Errorf("write login response error: %w", err)
	}

	connection.SetAuthenticated(true)
	connection.SetClientId(creds.ClientID)

	e.Metrics.Inc(ctx, metrics.AuthSuccessTotal)
	e.Log.WithSessionId(connection.SessionId()).WithUserId(creds.ClientID).Info("login successful")

	return nil
}

func handleLogout(ctx context.Context, connection *conn2.Connection, e *Epp) error {
	if connection.IsAuthenticated() {
		res := response.NewResponse[struct{}, struct{}](
			1500,
			response.CommandCompleteSuccessfullyEndingSession,
		)
		if err := connection.Write(ctx, res, e.Metrics.IncBytes); err != nil {
			return fmt.Errorf("write logout response error: %w", err)
		}

		return errSessionTerminate
	}

	errorResponse := response.AnyError(2002, response.CommandUseError)
	if err := connection.Write(ctx, errorResponse, e.Metrics.IncBytes); err != nil {
		return fmt.Errorf("write error response for invalid logout command: %w", err)
	}

	return nil
}
