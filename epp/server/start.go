package server

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"

	command2 "github.com/pixel365/zoner/epp/server/command"
	"github.com/pixel365/zoner/epp/server/command/command"
	login2 "github.com/pixel365/zoner/epp/server/command/login"
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
	defer func() {
		if err := connection.Close(); err != nil {
			e.Log.Error("close connection error", err)
		}
	}()

	greeting := command2.NewGreeting(e.Config.Greeting)
	if err := connection.Write(ctx, greeting); err != nil {
		e.Log.Error("write greeting error", err)
		return
	}

	parser := command2.CmdParser{}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			frame, err := connection.ReadFrame(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}

				e.Log.Error("read frame error", err)
				return
			}

			cmd, err := parseFrame(ctx, connection, &parser, frame)
			if err != nil {
				e.Log.Error("parse command error", err)
				return
			}

			if err = sendResponse(ctx, connection, cmd, e); err != nil {
				e.Log.Error("write frame error", err)
				return
			}
		}
	}
}

func parseFrame(
	ctx context.Context,
	connection *conn2.Connection,
	parser *command2.CmdParser,
	frame []byte,
) (command.Commander, error) {
	cmd, err := parser.Parse(frame)
	if err != nil {
		errorResponse := response.AnyError(2001, response.CommandSyntaxError)
		if err = connection.Write(ctx, errorResponse); err != nil {
			return nil, err
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
		greeting := command2.NewGreeting(e.Config.Greeting)
		if err := connection.Write(ctx, greeting); err != nil {
			return err
		}
		return nil
	}

	if cmd.Name() == command.Login {
		return handleLogin(ctx, connection, cmd, e)
	}

	if cmd.NeedAuth() && !connection.IsAuthenticated() {
		errorResponse := response.AnyError(2200, response.AuthorizationError)
		if err := connection.Write(ctx, errorResponse); err != nil {
			return err
		}
		return nil
	}

	errorResponse := response.AnyError(2101, response.UnimplementedCommand)
	if err := connection.Write(ctx, errorResponse); err != nil {
		return err
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
			return err
		}
		return nil
	}

	creds, ok := cmd.(login2.Login)
	if !ok {
		errorResponse := response.AnyError(2002, response.CommandUseError)
		if err := connection.Write(ctx, errorResponse); err != nil {
			return err
		}
		return nil
	}

	if err := e.AuthRepository.Login(creds.ClientID, creds.Password); err != nil {
		errorResponse := response.AnyError(2201, response.AuthorizationError)
		if err = connection.Write(ctx, errorResponse); err != nil {
			return err
		}

		return nil
	}

	connection.SetAuthenticated(true)
	res := response.NewResponse[struct{}, struct{}](1000, response.CommandCompletedSuccessfully)

	return connection.Write(ctx, res)
}
