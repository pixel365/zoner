package transfer

import (
	"encoding/xml"
	"errors"

	"github.com/pixel365/zoner/epp/server/command/internal"
)

func (t *Transfer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	t.Op = ""
	for _, a := range start.Attr {
		if a.Name.Local == "op" && a.Name.Space == "" {
			t.Op = a.Value
			break
		}
	}

	t.Domain, t.Contact = nil, nil

	var seen bool
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}

		done, err := t.handleToken(d, tok, &start, &seen)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}
}

func (t *Transfer) handleToken(
	d *xml.Decoder,
	tok xml.Token,
	start *xml.StartElement,
	seen *bool,
) (bool, error) {
	switch token := tok.(type) {
	case xml.StartElement:
		if token.Name.Local != t.Name().String() {
			return false, errors.New("unexpected element inside <transfer>: " + token.Name.Local)
		}

		if *seen {
			return false, errors.New("exactly one transfer command must be present")
		}

		if err := t.decodeObjectTransfer(d, &token); err != nil {
			return false, err
		}
		*seen = true

		return false, nil

	case xml.EndElement:
		if token.Name == start.Name {
			if !*seen {
				return false, errors.New("exactly one transfer command must be present")
			}

			return true, nil
		}
	}

	return false, nil
}

func (t *Transfer) decodeObjectTransfer(d *xml.Decoder, token *xml.StartElement) error {
	switch token.Name.Space {
	case internal.NsDomain:
		var x domainTransferXML
		if err := d.DecodeElement(&x, token); err != nil {
			return err
		}

		t.Domain = &Domain{Name: x.Name, AuthInfo: x.AuthInfo}

		return nil

	case internal.NsContact:
		var x contactTransferXML
		if err := d.DecodeElement(&x, token); err != nil {
			return err
		}

		t.Contact = &Contact{ID: x.ID, AuthInfo: x.AuthInfo}

		return nil

	default:
		return errors.New("unsupported <transfer> object namespace: " + token.Name.Space)
	}
}
