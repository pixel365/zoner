package command

import (
	"strings"
	"time"

	"github.com/pixel365/zoner/epp/config/epp/greeting"
)

// Greeting https://datatracker.ietf.org/doc/html/rfc5730#section-2.4
type Greeting struct {
	greeting greeting.Greeting
}

func (g Greeting) Marshal() ([]byte, error) {
	var b strings.Builder

	now := time.Now().UTC().Format(time.RFC3339)

	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	b.WriteString(`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">`)
	b.WriteString(`<greeting>`)

	b.WriteString(`<svID>` + g.greeting.ServerID + `</svID>`)
	b.WriteString(`<svDate>` + now + `</svDate>`)

	b.WriteString(`<svcMenu>`)

	for _, v := range g.greeting.Versions {
		b.WriteString(`<version>` + v + `</version>`)
	}

	for _, lang := range g.greeting.Languages {
		b.WriteString(`<lang>` + lang.String() + `</lang>`)
	}

	for _, o := range g.greeting.ObjURI {
		b.WriteString(`<objURI>` + o.String() + `</objURI>`)
	}

	if len(g.greeting.Extensions) > 0 {
		b.WriteString(`<svcExtension>`)

		for _, ex := range g.greeting.Extensions {
			b.WriteString(`<extURI>` + ex.String() + `</extURI>`)
		}

		b.WriteString(`</svcExtension>`)
	}

	if g.greeting.Dcp != nil {
		g.greeting.Dcp.WriteXML(&b)
	}

	b.WriteString(`</svcMenu>`)

	b.WriteString(`</greeting>`)
	b.WriteString(`</epp>`)

	return []byte(b.String()), nil
}

func NewGreeting(greeting greeting.Greeting) Greeting {
	return Greeting{greeting}
}
