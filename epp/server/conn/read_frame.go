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
		return nil, fmt.Errorf("read frame header: %w", err)
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
		return nil, fmt.Errorf("read frame payload: %w", err)
	}

	return payload, nil
}

func (c *Connection) setReadDeadline(ctx context.Context) error {
	if ddl, ok := ctx.Deadline(); ok {
		_ = c.conn.SetReadDeadline(ddl)
		return nil
	}

	if c.readTimeout > 0 {
		_ = c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
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
