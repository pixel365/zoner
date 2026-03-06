package tests

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
)

func netAddr() string {
	return eppHost + ":" + eppPort
}

func writeEPPFrame(conn net.Conn, payload []byte) error {
	if len(payload) == 0 {
		return errors.New("payload is empty")
	}

	return writeFrame(conn, payload)
}

func writeFrame(conn net.Conn, frame []byte) error {

	total64 := uint64(len(frame)) + uint64(minFrameLength)
	if total64 > math.MaxUint32 {
		return fmt.Errorf("frame too large")
	}

	total := uint32(total64)
	if total < minFrameLength {
		return fmt.Errorf("invalid frame length: %d", total)
	}

	if total > maxFrameSize {
		return fmt.Errorf("frame too large: %d (max %d)", total, maxFrameSize)
	}

	var h [4]byte
	binary.BigEndian.PutUint32(h[:], total)

	if err := writeAll(conn, h[:]); err != nil {
		return err
	}

	if len(frame) == 0 {
		return nil
	}

	return writeAll(conn, frame)
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

func readEPPFrame(r io.Reader) ([]byte, error) {
	var h [4]byte
	if _, err := io.ReadFull(r, h[:]); err != nil {
		return nil, fmt.Errorf("read frame header: %w", err)
	}

	total := binary.BigEndian.Uint32(h[:])
	if total < 4 {
		return nil, fmt.Errorf("invalid frame length: %d", total)
	}

	if total > maxFrameSize {
		return nil, fmt.Errorf("frame too large: %d (max %d)", total, maxFrameSize)
	}

	payloadLength := total - 4
	payload := make([]byte, payloadLength)

	if payloadLength == 0 {
		return payload, nil
	}

	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, fmt.Errorf("read frame payload: %w", err)
	}

	return payload, nil
}

func loginXML() string {
	payload := fmt.Sprintf(`
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <command>
        <login>
            <clID>%s</clID>
            <pw>%s</pw>
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
`, testingRegistrarUsername, testingRegistrarPassword)

	return payload
}
