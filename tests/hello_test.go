package tests

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/pixel365/goepp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHello(t *testing.T) {
	addr := netAddr()
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true, //nolint:gosec
	})

	require.NoError(t, err)

	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
	_, _ = readEPPFrame(conn)

	payload := `
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <hello/>
</epp>
`

	err = writeEPPFrame(conn, []byte(payload))
	require.NoError(t, err)

	resp, err := readEPPFrame(conn)
	require.NoError(t, err)

	parser := goepp.CmdParser{}
	_, err = parser.Parse(resp)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "exactly one command must be present")

	t.Cleanup(func() {
		_ = conn.Close()
	})
}
