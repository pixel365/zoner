package server

import "strconv"

//nolint:lll
var err2200 = []byte(
	`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><result code="2200"><msg>Authentication error</msg></result></response></epp>`,
)

var (
	eppPrefix = []byte(`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><result code="`)
	eppMid    = []byte(`"><msg>`)
	eppSuffix = []byte(`</msg></result></response></epp>`)
)

type Response struct {
	Msg  string
	Code int
}

func (e *Response) AsBytes() []byte {
	if e.Code == 2200 {
		return err2200
	}

	msg := e.Msg
	if msg == "" {
		msg = defaultMsg(e.Code)
	}

	b := make([]byte, 0, len(eppPrefix)+4+len(eppMid)+len(msg)+len(eppSuffix))
	b = append(b, eppPrefix...)
	b = strconv.AppendInt(b, int64(e.Code), 10)
	b = append(b, eppMid...)
	b = xmlEscapeAppend(b, msg)
	b = append(b, eppSuffix...)

	return b
}

//nolint:gocyclo,cyclop
func defaultMsg(code int) string {
	// see https://datatracker.ietf.org/doc/html/rfc5730#section-3
	switch code {
	case 1000:
		return "Command completed successfully"
	case 1001:
		return "Command completed successfully; action pending"
	case 1300:
		return "Command completed successfully; no messages"
	case 1301:
		return "Command completed successfully; ack to dequeue"
	case 1500:
		return "Command completed successfully; ending session"
	case 2000:
		return "Unknown command"
	case 2001:
		return "Command syntax error"
	case 2002:
		return "Command use error"
	case 2003:
		return "Required parameter missing"
	case 2004:
		return "Parameter value range error"
	case 2005:
		return "Parameter value syntax error"
	case 2100:
		return "Unimplemented protocol version"
	case 2101:
		return "Unimplemented command"
	case 2102:
		return "Unimplemented option"
	case 2103:
		return "Unimplemented extension"
	case 2104:
		return "Billing failure"
	case 2105:
		return "Object is not eligible for renewal"
	case 2106:
		return "Object is not eligible for transfer"
	case 2200:
		return "Authentication error"
	case 2201:
		return "Authorization error"
	case 2202:
		return "Invalid authorization information"
	case 2300:
		return "Object pending transfer"
	case 2301:
		return "Object not pending transfer"
	case 2302:
		return "Object exists"
	case 2303:
		return "Object does not exist"
	case 2304:
		return "Object status prohibits operation"
	case 2305:
		return "Object association prohibits operation"
	case 2306:
		return "Parameter value policy error"
	case 2307:
		return "Unimplemented object service"
	case 2308:
		return "Data management policy violation"
	case 2400:
		return "Command failed"
	case 2500:
		return "Command failed; server closing connection"
	case 2501:
		return "Authentication error; server closing connection"
	case 2502:
		return "Session limit exceeded; server closing connection"
	default:
		return "Unknown code"
	}
}

func xmlEscapeAppend(dst []byte, s string) []byte {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '&':
			dst = append(dst, "&amp;"...)
		case '<':
			dst = append(dst, "&lt;"...)
		case '>':
			dst = append(dst, "&gt;"...)
		case '"':
			dst = append(dst, "&quot;"...)
		case '\'':
			dst = append(dst, "&apos;"...)
		default:
			dst = append(dst, s[i])
		}
	}
	return dst
}
