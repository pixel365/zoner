package internal

import (
	"encoding/xml"
	"errors"
)

const (
	NsDomain  = "urn:ietf:params:xml:ns:domain-1.0"
	NsContact = "urn:ietf:params:xml:ns:contact-1.0"
	NsHost    = "urn:ietf:params:xml:ns:host-1.0"
)

type DomainContactHost interface {
	SetDomain(*Domain)
	GetDomain() *Domain
	SetContact(*Contact)
	GetContact() *Contact
	SetHost(*Host)
	GetHost() *Host
}

func HandleToken(
	dch DomainContactHost,
	d *xml.Decoder,
	tok xml.Token,
	start *xml.StartElement,
	seen *bool,
	name string,
) (bool, error) {
	switch t := tok.(type) {
	case xml.StartElement:
		if t.Name.Local != name {
			return false, errors.New("unexpected element inside <" + name + ">: " + t.Name.Local)
		}

		if *seen {
			return false, errors.New("exactly one info command must be present")
		}

		return false, DecodeObjectInfo(dch, d, &t, seen, name)
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

func DecodeObjectInfo(
	dch DomainContactHost,
	d *xml.Decoder,
	t *xml.StartElement,
	seen *bool,
	name string,
) error {
	switch t.Name.Space {
	case NsDomain:
		dch.SetDomain(new(Domain))
		return Decode[Domain](d, dch.GetDomain(), t, seen)
	case NsContact:
		dch.SetContact(new(Contact))
		return Decode[Contact](d, dch.GetContact(), t, seen)
	case NsHost:
		dch.SetHost(new(Host))
		return Decode[Host](d, dch.GetHost(), t, seen)
	default:
		return errors.New("unsupported <" + name + "> object namespace: " + t.Name.Space)
	}
}

func Decode[T any](d *xml.Decoder, to *T, t *xml.StartElement, seen *bool) error {
	var v T
	if err := d.DecodeElement(&v, t); err != nil {
		return err
	}

	*to = v
	*seen = true

	return nil
}
