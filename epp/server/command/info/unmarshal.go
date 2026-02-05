package info

import (
	"encoding/xml"
	"errors"
)

const (
	nsDomain  = "urn:ietf:params:xml:ns:domain-1.0"
	nsContact = "urn:ietf:params:xml:ns:contact-1.0"
	nsHost    = "urn:ietf:params:xml:ns:host-1.0"
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

		done, err := i.handleToken(d, tok, &start, &seen)
		if err != nil {
			return err
		}

		if done {
			return nil
		}
	}
}

func (i *Info) handleToken(
	d *xml.Decoder,
	tok xml.Token,
	start *xml.StartElement,
	seen *bool,
) (bool, error) {
	switch t := tok.(type) {
	case xml.StartElement:
		if t.Name.Local != i.Name().String() {
			return false, errors.New("unexpected element inside <info>: " + t.Name.Local)
		}

		if *seen {
			return false, errors.New("exactly one info command must be present")
		}

		return false, i.decodeObjectInfo(d, &t, seen)
	case xml.EndElement:
		if t.Name == start.Name {
			if !*seen {
				return false, errors.New("exactly one info command must be present")
			}

			return true, nil
		}
	}

	return false, nil
}

func (i *Info) decodeObjectInfo(d *xml.Decoder, t *xml.StartElement, seen *bool) error {
	switch t.Name.Space {
	case nsDomain:
		i.Domain = new(Domain)
		return decode[Domain](d, i.Domain, t, seen)
	case nsContact:
		i.Contact = new(Contact)
		return decode[Contact](d, i.Contact, t, seen)
	case nsHost:
		i.Host = new(Host)
		return decode[Host](d, i.Host, t, seen)
	default:
		return errors.New("unsupported <info> object namespace: " + t.Name.Space)
	}
}

func decode[T any](d *xml.Decoder, to *T, t *xml.StartElement, seen *bool) error {
	var v T
	if err := d.DecodeElement(&v, t); err != nil {
		return err
	}

	*to = v
	*seen = true

	return nil
}
