package tests

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckContactAvailable(t *testing.T) {
	addr := netAddr()
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true, //nolint:gosec
	})

	require.NoError(t, err)

	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
	_, _ = readEPPFrame(conn)

	err = writeEPPFrame(conn, []byte(loginXML()))
	require.NoError(t, err)

	resp, err := readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Command completed successfully")

	payload := `
<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <command>
        <check>
            <contact:check xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
                <contact:id>contact-check-free</contact:id>
            </contact:check>
        </check>
        <clTRID>CONTACT-CHECK-FREE</clTRID>
    </command>
</epp>
`

	err = writeEPPFrame(conn, []byte(payload))
	require.NoError(t, err)

	resp, err = readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Command completed successfully")
	assert.Contains(t, string(resp), `avail="1"`)
	assert.Contains(t, string(resp), ">contact-check-free<")
	assert.NotContains(t, string(resp), "In Use")

	t.Cleanup(func() {
		_ = conn.Close()
	})
}

func TestCheckContactInUse(t *testing.T) {
	addr := netAddr()
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true, //nolint:gosec
	})

	require.NoError(t, err)

	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
	_, _ = readEPPFrame(conn)

	err = writeEPPFrame(conn, []byte(loginXML()))
	require.NoError(t, err)

	resp, err := readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Command completed successfully")

	createPayload := `
<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <command>
        <create>
            <contact:create xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
                <contact:id>contact-check-busy</contact:id>

                <contact:postalInfo type="int">
                    <contact:name>John Doe</contact:name>
                    <contact:org>Example Inc.</contact:org>
                    <contact:addr>
                        <contact:street>123 Example Dr.</contact:street>
                        <contact:city>Dulles</contact:city>
                        <contact:cc>US</contact:cc>
                    </contact:addr>
                </contact:postalInfo>

                <contact:email>jdoe@example.com</contact:email>

                <contact:authInfo>
                    <contact:pw>password</contact:pw>
                </contact:authInfo>
            </contact:create>
        </create>
        <clTRID>CONTACT-CREATE-BUSY</clTRID>
    </command>
</epp>
`

	err = writeEPPFrame(conn, []byte(createPayload))
	require.NoError(t, err)

	resp, err = readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Command completed successfully")

	checkPayload := `
<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <command>
        <check>
            <contact:check xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
                <contact:id>contact-check-busy</contact:id>
            </contact:check>
        </check>
        <clTRID>CONTACT-CHECK-BUSY</clTRID>
    </command>
</epp>
`

	err = writeEPPFrame(conn, []byte(checkPayload))
	require.NoError(t, err)

	resp, err = readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Command completed successfully")
	assert.Contains(t, string(resp), `avail="0"`)
	assert.Contains(t, string(resp), ">contact-check-busy<")
	assert.Contains(t, string(resp), "In Use")

	t.Cleanup(func() {
		_ = conn.Close()
	})
}

func TestCheckContactMixedAvailability(t *testing.T) {
	addr := netAddr()
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true, //nolint:gosec
	})

	require.NoError(t, err)

	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
	_, _ = readEPPFrame(conn)

	err = writeEPPFrame(conn, []byte(loginXML()))
	require.NoError(t, err)

	resp, err := readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Command completed successfully")

	createPayload := `
<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <command>
        <create>
            <contact:create xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
                <contact:id>contact-check-mixed-busy</contact:id>
                <contact:postalInfo type="int">
                    <contact:name>John Doe</contact:name>
                    <contact:addr>
                        <contact:city>Dulles</contact:city>
                        <contact:cc>US</contact:cc>
                    </contact:addr>
                </contact:postalInfo>
                <contact:email>jdoe@example.com</contact:email>
                <contact:authInfo>
                    <contact:pw>password</contact:pw>
                </contact:authInfo>
            </contact:create>
        </create>
        <clTRID>CONTACT-CREATE-MIXED</clTRID>
    </command>
</epp>
`

	err = writeEPPFrame(conn, []byte(createPayload))
	require.NoError(t, err)

	resp, err = readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Command completed successfully")

	checkPayload := `
<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <command>
        <check>
            <contact:check xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
                <contact:id>contact-check-mixed-busy</contact:id>
                <contact:id>contact-check-mixed-free</contact:id>
            </contact:check>
        </check>
        <clTRID>CONTACT-CHECK-MIXED</clTRID>
    </command>
</epp>
`

	err = writeEPPFrame(conn, []byte(checkPayload))
	require.NoError(t, err)

	resp, err = readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Command completed successfully")
	assert.Contains(t, string(resp), `<contact:id avail="0">contact-check-mixed-busy</contact:id>`)
	assert.Contains(t, string(resp), `<contact:reason>In Use</contact:reason>`)
	assert.Contains(t, string(resp), `<contact:id avail="1">contact-check-mixed-free</contact:id>`)

	t.Cleanup(func() {
		_ = conn.Close()
	})
}
