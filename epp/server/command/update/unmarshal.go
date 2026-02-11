package update

import (
	"encoding/xml"
	"errors"

	"github.com/pixel365/zoner/epp/server/command/internal"
)

func (u *Update) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	u.Domain = nil

	var seen bool
	var done bool

	for !done {
		tok, err := d.Token()
		if err != nil {
			return err
		}

		if err := u.handleToken(d, tok, &start, &seen, &done); err != nil {
			return err
		}
	}

	return nil
}

func (u *Update) handleToken(
	d *xml.Decoder,
	tok xml.Token,
	start *xml.StartElement,
	seen *bool,
	done *bool,
) error {
	switch t := tok.(type) {
	case xml.StartElement:
		if t.Name.Local != u.Name().String() {
			return errors.New("unexpected element inside <update>: " + t.Name.Local)
		}

		if *seen {
			return errors.New("exactly one update object must be present")
		}

		if t.Name.Space != internal.NsDomain {
			return errors.New("unsupported <update> object namespace: " + t.Name.Space)
		}

		var x Domain
		if err := d.DecodeElement(&x, &t); err != nil {
			return err
		}

		u.Domain = &x
		*seen = true

		return nil

	case xml.EndElement:
		if t.Name == start.Name {
			if !*seen {
				return errors.New("exactly one update object must be present")
			}

			*done = true

			return nil
		}
	}

	return nil
}
