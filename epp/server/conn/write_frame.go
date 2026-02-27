package conn

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"time"

	"github.com/pixel365/zoner/epp/server/response"
	"github.com/pixel365/zoner/internal/observability/metrics"
)

func (c *Connection) Write(
	ctx context.Context,
	m response.Marshaller,
	fn metrics.IncBytesFunc,
) error {
	payload, err := m.Marshal()
	if err != nil {
		return err
	}

	if len(payload) == 0 {
		return errors.New("payload is empty")
	}

	if err = c.writeFrame(ctx, payload); err != nil {
		return err
	}

	fn(ctx, metrics.FramesWriteTotal, int64(len(payload)))

	return nil
}

func (c *Connection) writeFrame(ctx context.Context, frame []byte) error {
	if err := c.setWriteDeadline(ctx); err != nil {
		return err
	}
	defer c.resetWriteDeadline()

	total64 := uint64(len(frame)) + uint64(minFrameLength)
	if total64 > math.MaxUint32 {
		return fmt.Errorf("frame too large")
	}

	total := uint32(total64)
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
	if ddl, ok := ctx.Deadline(); ok {
		_ = c.conn.SetWriteDeadline(ddl)
		return nil
	}

	if c.writeTimeout > 0 {
		_ = c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
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
