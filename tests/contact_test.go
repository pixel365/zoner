package tests

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateContactSuccess(t *testing.T) {
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
<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <command>
        <create>
            <contact:create xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
                <contact:id>identifier</contact:id>

                <contact:postalInfo type="int">
                    <contact:name>John Doe</contact:name>
                    <contact:org>Example Inc.</contact:org>
                    <contact:addr>
                        <contact:street>123 Example Dr.</contact:street>
                        <contact:street>Suite 100</contact:street>
                        <contact:city>Dulles</contact:city>
                        <contact:sp>VA</contact:sp>
                        <contact:pc>20166-6503</contact:pc>
                        <contact:cc>US</contact:cc>
                    </contact:addr>
                </contact:postalInfo>

                <contact:postalInfo type="loc">
                    <contact:name>John Doe</contact:name>
                    <contact:org>Example LLC.</contact:org>
                    <contact:addr>
                        <contact:street>Example str.</contact:street>
                        <contact:city>NY</contact:city>
                        <contact:cc>US</contact:cc>
                    </contact:addr>
                </contact:postalInfo>

                <contact:voice x="123">+1.7035555555</contact:voice>
                <contact:fax>+1.7035555556</contact:fax>
                <contact:email>jdoe@example.com</contact:email>

                <contact:authInfo>
                    <contact:pw>password</contact:pw>
                </contact:authInfo>

                <contact:disclose flag="0">
                    <contact:voice/>
                    <contact:email/>
                    <contact:addr type="int"/>
                </contact:disclose>
            </contact:create>
        </create>
        <clTRID>CONTACT-CREATE-VALID-FULL</clTRID>
    </command>
</epp>
`

	err = writeEPPFrame(conn, []byte(payload))
	require.NoError(t, err)

	resp, err = readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Command completed successfully")

	t.Cleanup(func() {
		_ = conn.Close()
	})
}

func TestCreateInvalidScheme(t *testing.T) {
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
<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <command>
        <create>
            <contact:create xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
            </contact:create>
        </create>
        <clTRID>CONTACT-CREATE-VALID-FULL</clTRID>
    </command>
</epp>
`

	err = writeEPPFrame(conn, []byte(payload))
	require.NoError(t, err)

	resp, err = readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Command syntax error")

	t.Cleanup(func() {
		_ = conn.Close()
	})
}
