package conn

import (
	"net"
	"time"

	"github.com/pixel365/zoner/epp/config"
)

const minFrameLength = 4

type Connection struct {
	conn          net.Conn
	readTimeout   time.Duration
	writeTimeout  time.Duration
	idleTimeout   time.Duration
	maxFrameSize  uint32
	authenticated bool
}

func (c *Connection) SetAuthenticated(authenticated bool) {
	c.authenticated = authenticated
}

func (c *Connection) IsAuthenticated() bool {
	return c.authenticated
}

func (c *Connection) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func NewConnection(conn net.Conn, cfg *config.Epp) *Connection {
	return &Connection{
		conn:          conn,
		maxFrameSize:  4 * 1024 * 1024, // 4MB
		readTimeout:   time.Duration(cfg.ReadTimeout) * time.Second,
		writeTimeout:  time.Duration(cfg.WriteTimeout) * time.Second,
		idleTimeout:   time.Duration(cfg.IdleTimeout) * time.Second,
		authenticated: false,
	}
}
