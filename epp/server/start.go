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

	conn2 "github.com/pixel365/zoner/epp/server/conn"
	"github.com/pixel365/zoner/epp/server/response"
)

func (e *Epp) Start(ctx context.Context) error {
	ln, err := tls.Listen("tcp", e.Config.ListenAddr, e.TLSConfig)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		if err := ln.Close(); err != nil {
			e.Log.Error("close listener error", err)
		}
	}()

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

func (e *Epp) handleConnection(ctx context.Context, conn net.Conn) {
	connection := conn2.NewConnection(conn, &e.Config)
	log := e.Log.SessionId(connection.SessionId())

	log.Info("new connection")

	defer func() {
		if err := connection.Close(); err != nil {
			log.ClientId(connection.ClientId()).Error("close connection error", err)
			return
		}

		duration := time.Since(connection.SessionStart())
		log.ClientId(connection.ClientId()).
			Info("connection closed. session duration: %s", duration.String())
	}()

	g := greeting.NewGreeting(e.Config.Greeting)
	if err := connection.Write(ctx, g); err != nil {
		log.Error("write greeting error", err)
		return
	}

	parser := goepp.CmdParser{}

	for {
		select {
		case <-ctx.Done():
			duration := time.Since(connection.SessionStart())
			log.ClientId(connection.ClientId()).
				Info("context done. session duration: %s", duration.String())
			return
		default:
			clientId := connection.ClientId()
			frame, err := connection.ReadFrame(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}

				log.ClientId(clientId).Error("read frame error", err)
				return
			}

			cmd, err := parseFrame(ctx, connection, &parser, frame)
			if err != nil {
				log.ClientId(clientId).Error("parse command error", err)
				return
			}

			if err = sendResponse(ctx, connection, cmd, e); err != nil {
				log.ClientId(clientId).Error("write frame error", err)
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
) (command.Commander, error) {
	cmd, err := parser.Parse(frame)
	if err != nil {
		errorResponse := response.AnyError(2001, response.CommandSyntaxError)
		if err = connection.Write(ctx, errorResponse); err != nil {
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
		if err := connection.Write(ctx, g); err != nil {
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
		if err := connection.Write(ctx, errorResponse); err != nil {
			return fmt.Errorf("write error response when client is not authenticated: %w", err)
		}
		return nil
	}

	errorResponse := response.AnyError(2101, response.UnimplementedCommand)
	if err := connection.Write(ctx, errorResponse); err != nil {
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
		if err := connection.Write(ctx, errorResponse); err != nil {
			return fmt.Errorf("write error response when client is authenticated: %w", err)
		}
		return nil
	}

	creds, ok := cmd.(login.Login)
	if !ok {
		errorResponse := response.AnyError(2002, response.CommandUseError)
		if err := connection.Write(ctx, errorResponse); err != nil {
			return fmt.Errorf("write error response for invalid login command: %w", err)
		}
		return nil
	}

	if err := e.AuthRepository.Login(creds.ClientID, creds.Password); err != nil {
		errorResponse := response.AnyError(2201, response.AuthorizationError)
		if err = connection.Write(ctx, errorResponse); err != nil {
			return fmt.Errorf("write error response for invalid login credentials: %w", err)
		}

		return nil
	}

	res := response.NewResponse[struct{}, struct{}](1000, response.CommandCompletedSuccessfully)
	if err := connection.Write(ctx, res); err != nil {
		return fmt.Errorf("write login response error: %w", err)
	}

	connection.SetAuthenticated(true)
	connection.SetClientId(creds.ClientID)

	e.Log.SessionId(connection.SessionId()).ClientId(creds.ClientID).Info("login successful")

	return nil
}

func handleLogout(ctx context.Context, connection *conn2.Connection, e *Epp) error {
	if connection.IsAuthenticated() {
		res := response.NewResponse[struct{}, struct{}](
			1500,
			response.CommandCompleteSuccessfullyEndingSession,
		)
		if err := connection.Write(ctx, res); err != nil {
			return fmt.Errorf("write logout response error: %w", err)
		}

		if err := connection.Close(); err != nil {
			return fmt.Errorf("close connection error: %w", err)
		}

		e.Log.SessionId(connection.SessionId()).
			ClientId(connection.ClientId()).
			Info("logout successful")

		return nil
	}

	errorResponse := response.AnyError(2002, response.CommandUseError)
	if err := connection.Write(ctx, errorResponse); err != nil {
		return fmt.Errorf("write error response for invalid logout command: %w", err)
	}

	return nil
}
