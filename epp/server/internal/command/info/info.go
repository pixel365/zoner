package info

import (
	"errors"
)

type Info struct {
	Domain  *Domain  `xml:"urn:ietf:params:xml:ns:domain-1.0 info"`
	Contact *Contact `xml:"urn:ietf:params:xml:ns:contact-1.0 info"`
	Host    *Host    `xml:"urn:ietf:params:xml:ns:host-1.0 info"`
}

func (i *Info) Name() string {
	return "info"
}

func (i *Info) Validate() error {
	var notNil uint8

	if i.Domain != nil {
		notNil++
	}

	if i.Contact != nil {
		notNil++
	}

	if i.Host != nil {
		notNil++
	}

	if notNil != 1 {
		return errors.New("exactly one info command must be present")
	}

	return i.validate()
}

func (i *Info) validate() error {
	if i.Domain != nil {
		return i.Domain.Validate()
	}

	if i.Contact != nil {
		return i.Contact.Validate()
	}

	if i.Host != nil {
		return i.Host.Validate()
	}

	return nil
}
