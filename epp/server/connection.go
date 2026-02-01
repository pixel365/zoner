package server

import (
	"context"
	"net"
)

type Connection struct {
	conn net.Conn
}

func (c *Connection) ReadFrame(_ context.Context) ([]byte, error) {
	return nil, nil
}

func (c *Connection) WriteFrame(_ context.Context, frame []byte) error {
	return nil
}

func (c *Connection) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{conn: conn}
}
