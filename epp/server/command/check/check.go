package check

import (
	"errors"

	"github.com/pixel365/zoner/epp/server/command/command"
)

type Check struct {
	Domain  *DomainCheck  `xml:"urn:ietf:params:xml:ns:domain-1.0"`
	Contact *ContactCheck `xml:"urn:ietf:params:xml:ns:contact-1.0"`
	Host    *HostCheck    `xml:"urn:ietf:params:xml:ns:host-1.0"`
}

func (c *Check) Name() command.CommandName {
	return command.Check
}

func (c *Check) NeedAuth() bool {
	return true
}

func (c *Check) validate() error {
	switch {
	case c.Domain != nil:
		return c.validateDomains()

	case c.Contact != nil:
		return c.validateContacts()

	case c.Host != nil:
		return c.validateHosts()
	}

	return nil
}

func (c *Check) validateDomains() error {
	if len(c.Domain.Names) == 0 {
		return errors.New("domain names is empty")
	}

	for _, n := range c.Domain.Names {
		if n == "" {
			return errors.New("domain name is empty")
		}
	}

	return nil
}

func (c *Check) validateContacts() error {
	if len(c.Contact.IDs) == 0 {
		return errors.New("contact ids is empty")
	}

	for _, id := range c.Contact.IDs {
		if id == "" {
			return errors.New("contact id is empty")
		}
	}

	return nil
}

func (c *Check) validateHosts() error {
	if len(c.Host.Names) == 0 {
		return errors.New("host names is empty")
	}

	for _, n := range c.Host.Names {
		if n == "" {
			return errors.New("host name is empty")
		}
	}

	return nil
}

func (c *Check) Validate() error {
	var notNil uint8

	if c.Domain != nil {
		notNil++
	}

	if c.Contact != nil {
		notNil++
	}

	if c.Host != nil {
		notNil++
	}

	if notNil != 1 {
		return errors.New("exactly one check command must be present")
	}

	return c.validate()
}
