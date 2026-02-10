package renew

import (
	"encoding/xml"
	"errors"

	"github.com/pixel365/zoner/epp/server/command/internal"
)

func (r *Renew) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	r.Domain = nil

	var seen bool
	var done bool

	for !done {
		tok, err := d.Token()
		if err != nil {
			return err
		}

		if err := r.handleToken(d, tok, &start, &seen, &done); err != nil {
			return err
		}
	}

	return nil
}

func (r *Renew) handleToken(
	d *xml.Decoder,
	tok xml.Token,
	start *xml.StartElement,
	seen *bool,
	done *bool,
) error {
	switch t := tok.(type) {
	case xml.StartElement:
		if t.Name.Local != r.Name().String() {
			return errors.New("unexpected element inside <renew>: " + t.Name.Local)
		}

		if *seen {
			return errors.New("exactly one <renew> object must be present")
		}

		if err := r.decodeObject(d, &t); err != nil {
			return err
		}

		*seen = true

		return nil

	case xml.EndElement:
		if t.Name == start.Name {
			if !*seen {
				return errors.New("exactly one <renew> object must be present")
			}

			*done = true

			return nil
		}
	}

	return nil
}

func (r *Renew) decodeObject(d *xml.Decoder, t *xml.StartElement) error {
	if t.Name.Local != r.Name().String() {
		return errors.New("unexpected element inside <renew>: " + t.Name.Local)
	}

	if t.Name.Space != internal.NsDomain {
		return errors.New("unsupported <renew> object namespace: " + t.Name.Space)
	}

	var x Domain
	if err := d.DecodeElement(&x, t); err != nil {
		return err
	}

	r.Domain = &x

	return nil
}
