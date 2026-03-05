package tests

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/pixel365/goepp"
	"github.com/stretchr/testify/require"
)

func TestGreetingIncoming(t *testing.T) {
	addr := netAddr()
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true, //nolint:gosec
	})

	require.NoError(t, err)

	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))

	greeting, err := readEPPFrame(conn)
	parser := goepp.CmdParser{}
	_, err = parser.Parse(greeting)

	require.Error(t, err)
	require.Equal(t, "exactly one command must be present", err.Error())
}
