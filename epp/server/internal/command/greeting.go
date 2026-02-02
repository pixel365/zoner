package command

import (
	"strings"
	"time"

	"github.com/pixel365/zoner/epp/config/epp/greeting"
)

// Greeting https://datatracker.ietf.org/doc/html/rfc5730#section-2.4
type Greeting struct{}

func (g Greeting) Bytes(greeting greeting.Greeting) []byte {
	var b strings.Builder

	now := time.Now().UTC().Format(time.RFC3339)

	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	b.WriteString(`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">`)
	b.WriteString(`<greeting>`)

	b.WriteString(`<svID>` + greeting.ServerID + `</svID>`)
	b.WriteString(`<svDate>` + now + `</svDate>`)

	b.WriteString(`<svcMenu>`)

	for _, v := range greeting.Versions {
		b.WriteString(`<version>` + v + `</version>`)
	}

	for _, lang := range greeting.Languages {
		b.WriteString(`<lang>` + lang.String() + `</lang>`)
	}

	for _, o := range greeting.ObjURI {
		b.WriteString(`<objURI>` + o.String() + `</objURI>`)
	}

	if len(greeting.Extensions) > 0 {
		b.WriteString(`<svcExtension>`)

		for _, ex := range greeting.Extensions {
			b.WriteString(`<extURI>` + ex.String() + `</extURI>`)
		}

		b.WriteString(`</svcExtension>`)
	}

	if greeting.Dcp != nil {
		greeting.Dcp.WriteXML(&b)
	}

	b.WriteString(`</svcMenu>`)

	b.WriteString(`</greeting>`)
	b.WriteString(`</epp>`)

	return []byte(b.String())
}

func NewGreeting() Greeting {
	return Greeting{}
}
