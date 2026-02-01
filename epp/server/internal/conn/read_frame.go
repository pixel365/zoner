package conn

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

func (c *Connection) ReadFrame(ctx context.Context) ([]byte, error) {
	if err := c.setReadDeadline(ctx); err != nil {
		return nil, err
	}
	defer c.resetReadDeadline()

	var h [4]byte
	if _, err := io.ReadFull(c.conn, h[:]); err != nil {
		return nil, err
	}

	total := binary.BigEndian.Uint32(h[:])
	if total < minFrameLength {
		return nil, fmt.Errorf("invalid frame length: %d", total)
	}

	if c.maxFrameSize > 0 && total > c.maxFrameSize {
		return nil, fmt.Errorf("frame too large: %d (max %d)", total, c.maxFrameSize)
	}

	payloadLength := total - minFrameLength
	payload := make([]byte, payloadLength)

	if payloadLength == 0 {
		return payload, nil
	}

	if _, err := io.ReadFull(c.conn, payload); err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Connection) setReadDeadline(ctx context.Context) error {
	if c.readTimeout > 0 {
		_ = c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	}

	if ddl, ok := ctx.Deadline(); ok {
		_ = c.conn.SetReadDeadline(ddl)
		return nil
	}

	if ctx.Done() != nil {
		_ = c.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	}

	return nil
}

func (c *Connection) resetReadDeadline() {
	_ = c.conn.SetReadDeadline(time.Time{})
}
