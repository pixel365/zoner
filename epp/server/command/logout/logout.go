package logout

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/pixel365/zoner/epp/server/command/command"
)

type Logout struct{}

func (l *Logout) Name() command.CommandName {
	return command.Logout
}

func (l *Logout) NeedAuth() bool {
	return true
}

func (l *Logout) Validate() error {
	return nil
}

func (l *Logout) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
				return fmt.Errorf("<logout> must be empty")
			}
		default:
			return fmt.Errorf("<logout> must be empty")
		}
	}
}
