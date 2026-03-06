package tests

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogoutSuccess(t *testing.T) {
	addr := netAddr()
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true, //nolint:gosec
	})

	require.NoError(t, err)

	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
	_, _ = readEPPFrame(conn)

	payload := loginXML()

	err = writeEPPFrame(conn, []byte(payload))
	require.NoError(t, err)

	resp, err := readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Command completed successfully")

	payload = `
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <command>
        <logout/>
        <clTRID>ABC-12345</clTRID>
    </command>
</epp>
`

	err = writeEPPFrame(conn, []byte(payload))
	require.NoError(t, err)

	resp, err = readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Command completed successfully; ending session")

	t.Cleanup(func() {
		_ = conn.Close()
	})
}

func TestUnauthorizedLogoutFail(t *testing.T) {
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
    <command>
        <logout/>
        <clTRID>ABC-12345</clTRID>
    </command>
</epp>
`

	err = writeEPPFrame(conn, []byte(payload))
	require.NoError(t, err)

	resp, err := readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Authorization error")

	t.Cleanup(func() {
		_ = conn.Close()
	})
}
