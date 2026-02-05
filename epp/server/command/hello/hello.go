package hello

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/pixel365/zoner/epp/server/command/command"
)

type Hello struct{}

func (h *Hello) Name() command.CommandName {
	return command.Hello
}

func (h *Hello) NeedAuth() bool {
	return false
}

func (h *Hello) Validate() error {
	return nil
}

func (h *Hello) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.EndElement:
			if t.Name == start.Name {
				return nil
			}
		case xml.CharData:
			if len(bytes.TrimSpace(t)) != 0 {
				return fmt.Errorf("<hello> must be empty")
			}
		default:
			return fmt.Errorf("<hello> must be empty")
		}
	}
}
