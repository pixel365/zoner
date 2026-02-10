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

type DomainContactHost[D any, C any, H any] interface {
	SetDomain(*D)
	GetDomain() *D
	SetContact(*C)
	GetContact() *C
	SetHost(*H)
	GetHost() *H
}

func HandleToken[D any, C any, H any](
	dch DomainContactHost[D, C, H],
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

		if err := DecodeObjectInfo(dch, d, &t, name); err != nil {
			return false, err
		}

		*seen = true

		return false, nil
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

func DecodeObjectInfo[D any, C any, H any](
	dch DomainContactHost[D, C, H],
	d *xml.Decoder,
	t *xml.StartElement,
	name string,
) error {
	switch t.Name.Space {
	case NsDomain:
		dch.SetDomain(new(D))
		return Decode[D](d, dch.GetDomain(), t)
	case NsContact:
		dch.SetContact(new(C))
		return Decode[C](d, dch.GetContact(), t)
	case NsHost:
		dch.SetHost(new(H))
		return Decode[H](d, dch.GetHost(), t)
	default:
		return errors.New("unsupported <" + name + "> object namespace: " + t.Name.Space)
	}
}

func Decode[T any](d *xml.Decoder, to *T, t *xml.StartElement) error {
	return d.DecodeElement(to, t)
}
