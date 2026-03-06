package tests

import (
	"crypto/tls"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginInvalidCredentials(t *testing.T) {
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
        <login>
            <clID>ClientX</clID>
            <pw>foo-BAR2</pw>
            <newPW>bar-FOO2</newPW>
            <options>
                <version>1.0</version>
                <lang>en</lang>
            </options>
            <svcs>
                <objURI>urn:ietf:params:xml:ns:obj1</objURI>
                <objURI>urn:ietf:params:xml:ns:obj2</objURI>
                <objURI>urn:ietf:params:xml:ns:obj3</objURI>
                <svcExtension>
                    <extURI>http://custom/obj1ext-1.0</extURI>
                </svcExtension>
            </svcs>
        </login>
        <clTRID>ABC-12345</clTRID>
    </command>
</epp>
`

	err = writeEPPFrame(conn, []byte(payload))
	require.NoError(t, err)

	resp, err := readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Authentication error")

	t.Cleanup(func() {
		_ = conn.Close()
	})
}

func TestLoginValidCredentials(t *testing.T) {
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

	t.Cleanup(func() {
		_ = conn.Close()
	})
}

func TestAlreadyLoggedIn(t *testing.T) {
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

	payload = loginXML()

	err = writeEPPFrame(conn, []byte(payload))
	require.NoError(t, err)

	resp, err = readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Already logged in")

	t.Cleanup(func() {
		_ = conn.Close()
	})
}

func TestLoginChangePassword(t *testing.T) {
	addr := netAddr()
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true, //nolint:gosec
	})

	require.NoError(t, err)

	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
	_, _ = readEPPFrame(conn)

	oldPassword := testingRegistrarPassword

	payload := fmt.Sprintf(`
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <command>
        <login>
            <clID>%s</clID>
            <pw>%s</pw>
            <newPW>%s</newPW>
            <options>
                <version>1.0</version>
                <lang>en</lang>
            </options>
            <svcs>
                <objURI>urn:ietf:params:xml:ns:obj1</objURI>
                <objURI>urn:ietf:params:xml:ns:obj2</objURI>
                <objURI>urn:ietf:params:xml:ns:obj3</objURI>
                <svcExtension>
                    <extURI>http://custom/obj1ext-1.0</extURI>
                </svcExtension>
            </svcs>
        </login>
        <clTRID>ABC-12345</clTRID>
    </command>
</epp>
`, testingRegistrarUsername, testingRegistrarPassword, "12345678")

	err = writeEPPFrame(conn, []byte(payload))
	require.NoError(t, err)

	resp, err := readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Command completed successfully")

	testingRegistrarPassword = "12345678"

	payload = fmt.Sprintf(`
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <command>
        <login>
            <clID>%s</clID>
            <pw>%s</pw>
            <newPW>%s</newPW>
            <options>
                <version>1.0</version>
                <lang>en</lang>
            </options>
            <svcs>
                <objURI>urn:ietf:params:xml:ns:obj1</objURI>
                <objURI>urn:ietf:params:xml:ns:obj2</objURI>
                <objURI>urn:ietf:params:xml:ns:obj3</objURI>
                <svcExtension>
                    <extURI>http://custom/obj1ext-1.0</extURI>
                </svcExtension>
            </svcs>
        </login>
        <clTRID>ABC-12345</clTRID>
    </command>
</epp>
`, testingRegistrarUsername, testingRegistrarPassword, oldPassword)

	err = writeEPPFrame(conn, []byte(payload))
	require.NoError(t, err)

	resp, err = readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "Command completed successfully")

	testingRegistrarPassword = oldPassword

	t.Cleanup(func() {
		_ = conn.Close()
	})
}

func TestLoginChangePasswordFail(t *testing.T) {
	addr := netAddr()
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true, //nolint:gosec
	})

	require.NoError(t, err)

	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
	_, _ = readEPPFrame(conn)

	payload := fmt.Sprintf(`
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <command>
        <login>
            <clID>%s</clID>
            <pw>%s</pw>
            <newPW>%s</newPW>
            <options>
                <version>1.0</version>
                <lang>en</lang>
            </options>
            <svcs>
                <objURI>urn:ietf:params:xml:ns:obj1</objURI>
                <objURI>urn:ietf:params:xml:ns:obj2</objURI>
                <objURI>urn:ietf:params:xml:ns:obj3</objURI>
                <svcExtension>
                    <extURI>http://custom/obj1ext-1.0</extURI>
                </svcExtension>
            </svcs>
        </login>
        <clTRID>ABC-12345</clTRID>
    </command>
</epp>
`, testingRegistrarUsername, testingRegistrarPassword, "123")

	err = writeEPPFrame(conn, []byte(payload))
	require.NoError(t, err)

	resp, err := readEPPFrame(conn)
	require.NoError(t, err)

	assert.Contains(t, string(resp), "password length must be at least")

	t.Cleanup(func() {
		_ = conn.Close()
	})
}
