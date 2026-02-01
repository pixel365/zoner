package conn

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func (c *Connection) WriteFrame(ctx context.Context, frame []byte) error {
	if err := c.setWriteDeadline(ctx); err != nil {
		return err
	}
	defer c.resetWriteDeadline()

	//nolint:gosec
	total := uint32(len(frame)) + minFrameLength
	if total < minFrameLength {
		return fmt.Errorf("invalid frame length: %d", total)
	}

	if c.maxFrameSize > 0 && total > c.maxFrameSize {
		return fmt.Errorf("frame too large: %d (max %d)", total, c.maxFrameSize)
	}

	var h [4]byte
	binary.BigEndian.PutUint32(h[:], total)

	if err := writeAll(c.conn, h[:]); err != nil {
		return err
	}

	if len(frame) == 0 {
		return nil
	}

	return writeAll(c.conn, frame)
}

func (c *Connection) setWriteDeadline(ctx context.Context) error {
	if c.writeTimeout > 0 {
		_ = c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}

	if dl, ok := ctx.Deadline(); ok {
		_ = c.conn.SetWriteDeadline(dl)
		return nil
	}

	if ctx.Done() != nil {
		_ = c.conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	}

	return nil
}

func (c *Connection) resetWriteDeadline() {
	_ = c.conn.SetWriteDeadline(time.Time{})
}

func writeAll(conn net.Conn, b []byte) error {
	for len(b) > 0 {
		n, err := conn.Write(b)
		if err != nil {
			return err
		}
		b = b[n:]
	}

	return nil
}
