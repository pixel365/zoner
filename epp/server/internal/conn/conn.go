package conn

import (
	"net"
	"time"

	"github.com/pixel365/zoner/internal/stringutils"
)

const minFrameLength = 4

var (
	frameReadTtl  time.Duration
	frameWriteTtl time.Duration
	connIdleTtl   time.Duration
)

func init() {
	readTtl := stringutils.GetPositiveIntFromEnv("EPP_FRAME_READ_TTL", 30)
	writeTtl := stringutils.GetPositiveIntFromEnv("EPP_FRAME_WRITE_TTL", 10)
	idleTtl := stringutils.GetPositiveIntFromEnv("EPP_CONN_IDLE_TTL", 1800) // default 30 minutes

	frameReadTtl = time.Duration(readTtl) * time.Second
	frameWriteTtl = time.Duration(writeTtl) * time.Second
	connIdleTtl = time.Duration(idleTtl) * time.Second
}

type Connection struct {
	conn         net.Conn
	maxFrameSize uint32
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
}

func (c *Connection) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn:         conn,
		maxFrameSize: 4 * 1024 * 1024, // 4MB
		readTimeout:  frameReadTtl,
		writeTimeout: frameWriteTtl,
		idleTimeout:  connIdleTtl,
	}
}
