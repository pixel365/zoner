package create

import (
	"errors"

	"github.com/pixel365/zoner/epp/server/command/command"
)

type Create struct {
	Domain  *Domain
	Contact *Contact
	Host    *Host
}

func (c *Create) Name() command.CommandName {
	return command.Create
}

func (c *Create) NeedAuth() bool {
	return true
}

func (c *Create) Validate() error {
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
		return errors.New("exactly one create command must be present")
	}

	switch {
	case c.Domain != nil:
		return c.Domain.Validate()
	case c.Contact != nil:
		return c.Contact.Validate()
	case c.Host != nil:
		return c.Host.Validate()
	}

	return nil
}
