package info

import (
	"encoding/xml"

	"github.com/pixel365/zoner/epp/server/command/internal"
)

func (i *Info) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	i.Domain = nil
	i.Contact = nil
	i.Host = nil

	var seen bool

	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}

		done, err := internal.HandleToken(i, d, tok, &start, &seen, i.Name().String())
		if err != nil {
			return err
		}

		if done {
			return nil
		}
	}
}
