package server

import (
	"context"
	"crypto/tls"
	"net"

	conn2 "github.com/pixel365/zoner/epp/server/internal/conn"
)

func (e *Epp) Start(ctx context.Context) error {
	ln, err := tls.Listen("tcp", e.Addr, e.TLSConfig)
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
	connection := conn2.NewConnection(conn)
	defer func() {
		if err := connection.Close(); err != nil {
			e.Log.Error("close connection error", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			frame, err := connection.ReadFrame(ctx)
			if err != nil {
				e.Log.Error("read frame error", err)
				return
			}

			if err = connection.WriteFrame(ctx, frame); err != nil {
				e.Log.Error("write frame error", err)
				return
			}
		}
	}
}
