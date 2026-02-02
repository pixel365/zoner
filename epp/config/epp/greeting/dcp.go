package greeting

import (
	"errors"
	"strings"
)

type DcpAccess string

func (d DcpAccess) String() string {
	return string(d)
}

const (
	DcpAll              DcpAccess = "all"
	DcpNone             DcpAccess = "none"
	DcpNull             DcpAccess = "null"
	DcpPersonalAndOther DcpAccess = "personalAndOther"
	DcpOther            DcpAccess = "other"
)

type Dcp struct {
	Access     DcpAccess   `yaml:"access"`
	Statements []Statement `yaml:"statements"`
}

func (d Dcp) Validate() error {
	if d.Access == "" {
		return errors.New("access is empty")
	}

	if len(d.Statements) == 0 {
		return errors.New("statements is empty")
	}

	for _, s := range d.Statements {
		if err := s.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (d Dcp) WriteXML(b *strings.Builder) {
	b.WriteString(`<dcp>`)
	b.WriteString(`<access><` + d.Access.String() + `/></access>`)

	for _, s := range d.Statements {
		b.WriteString(`<statement>`)
		s.WriteXML(b)
		b.WriteString(`</statement>`)
	}

	b.WriteString(`</dcp>`)
}
