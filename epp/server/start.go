package server

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"

	"github.com/pixel365/zoner/epp/server/internal/command"
	conn2 "github.com/pixel365/zoner/epp/server/internal/conn"
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

	var greeting command.Greeting
	if err := connection.WriteFrame(ctx, greeting.Bytes(e.Config.Greeting)); err != nil {
		e.Log.Error("write greeting error", err)
		return
	}

	parser := command.CmdParser{}

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

			cmd, err := parser.Parse(frame)
			if err != nil {
				//TODO: send error response
				e.Log.Error("parse command error", err)
				continue
			}

			if err = sendResponse(ctx, connection, cmd, e); err != nil {
				e.Log.Error("write frame error", err)
				return
			}
		}
	}
}

func sendResponse(
	ctx context.Context,
	connection *conn2.Connection,
	cmd command.Command,
	e *Epp,
) error {
	if cmd.Name() == "hello" {
		var greeting command.Greeting
		if err := connection.WriteFrame(ctx, greeting.Bytes(e.Config.Greeting)); err != nil {
			return err
		}
		return nil
	}

	var response command.Response
	if err := connection.WriteFrame(ctx, response.AsBytes()); err != nil {
		return err
	}

	return nil
}
