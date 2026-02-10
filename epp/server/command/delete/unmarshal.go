package delete

import (
	"encoding/xml"

	"github.com/pixel365/zoner/epp/server/command/internal"
)

func (d *Delete) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	d.Domain = nil
	d.Contact = nil
	d.Host = nil

	var seen bool

	for {
		tok, err := dec.Token()
		if err != nil {
			return err
		}

		done, err := internal.HandleToken(d, dec, tok, &start, &seen, d.Name().String())
		if err != nil {
			return err
		}

		if done {
			return nil
		}
	}
}
