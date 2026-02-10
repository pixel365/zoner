package delete

import (
	"errors"

	"github.com/pixel365/zoner/epp/server/command/command"
)

type Delete struct {
	Domain  *DomainDelete  `xml:"urn:ietf:params:xml:ns:domain-1.0 delete"`
	Contact *ContactDelete `xml:"urn:ietf:params:xml:ns:contact-1.0 delete"`
	Host    *HostDelete    `xml:"urn:ietf:params:xml:ns:host-1.0 delete"`
}

func (d *Delete) Name() command.CommandName {
	return command.Delete
}

func (d *Delete) NeedAuth() bool {
	return true
}

func (d *Delete) Validate() error {
	var notNil uint8

	if d.Domain != nil {
		notNil++
	}

	if d.Contact != nil {
		notNil++
	}

	if d.Host != nil {
		notNil++
	}

	if notNil != 1 {
		return errors.New("exactly one delete command must be present")
	}

	switch {
	case d.Domain != nil:
		return d.Domain.Validate()
	case d.Contact != nil:
		return d.Contact.Validate()
	case d.Host != nil:
		return d.Host.Validate()
	}

	return nil
}
