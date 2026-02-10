package command

import (
	"encoding/xml"
	"errors"

	"github.com/pixel365/zoner/epp/server/command/check"
	command2 "github.com/pixel365/zoner/epp/server/command/command"
	"github.com/pixel365/zoner/epp/server/command/create"
	"github.com/pixel365/zoner/epp/server/command/delete"
	"github.com/pixel365/zoner/epp/server/command/hello"
	"github.com/pixel365/zoner/epp/server/command/info"
	"github.com/pixel365/zoner/epp/server/command/login"
	"github.com/pixel365/zoner/epp/server/command/logout"
	"github.com/pixel365/zoner/epp/server/command/poll"
	"github.com/pixel365/zoner/epp/server/command/renew"
	"github.com/pixel365/zoner/epp/server/command/transfer"
)

type EPP struct {
	// Hello https://datatracker.ietf.org/doc/html/rfc5730#section-2.3
	Hello *hello.Hello `xml:"urn:ietf:params:xml:ns:epp-1.0 hello"`

	// Login https://datatracker.ietf.org/doc/html/rfc5730#section-2.9.1.1
	Login *login.Login `xml:"urn:ietf:params:xml:ns:epp-1.0 command>login"`

	// Logout https://datatracker.ietf.org/doc/html/rfc5730#section-2.9.1.2
	Logout *logout.Logout `xml:"urn:ietf:params:xml:ns:epp-1.0 command>logout"`

	// Check https://datatracker.ietf.org/doc/html/rfc5730#section-2.9.2.1
	Check *check.Check `xml:"urn:ietf:params:xml:ns:epp-1.0 command>check"`

	// Info https://datatracker.ietf.org/doc/html/rfc5730#section-2.9.2.2
	Info *info.Info `xml:"urn:ietf:params:xml:ns:epp-1.0 command>info"`

	// Poll https://datatracker.ietf.org/doc/html/rfc5730#section-2.9.2.3
	Poll *poll.Poll `xml:"urn:ietf:params:xml:ns:epp-1.0 command>poll"`

	// Transfer https://datatracker.ietf.org/doc/html/rfc5730#section-2.9.2.4
	// https://datatracker.ietf.org/doc/html/rfc5730#section-2.9.3.4
	Transfer *transfer.Transfer `xml:"urn:ietf:params:xml:ns:epp-1.0 command>transfer"`

	// Create https://datatracker.ietf.org/doc/html/rfc5730#section-2.9.3.1
	Create *create.Create `xml:"urn:ietf:params:xml:ns:epp-1.0 command>create"`

	// Delete https://datatracker.ietf.org/doc/html/rfc5730#section-2.9.3.2
	Delete *delete.Delete `xml:"urn:ietf:params:xml:ns:epp-1.0 command>delete"`

	// Renew https://datatracker.ietf.org/doc/html/rfc5730#section-2.9.3.3
	Renew *renew.Renew `xml:"urn:ietf:params:xml:ns:epp-1.0 command>renew"`

	// XMLName https://datatracker.ietf.org/doc/html/rfc5730#section-2.2
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`

	TransactionID string `xml:"urn:ietf:params:xml:ns:epp-1.0 command>clTRID,omitempty"` //TODO: add validator
}

func (e *EPP) Validate() error { //TODO: change to error response
	var notNilCommands uint8

	if e.Hello != nil {
		notNilCommands++
	}

	if e.Login != nil {
		notNilCommands++
	}

	if e.Logout != nil {
		notNilCommands++
	}

	if e.Check != nil {
		notNilCommands++
	}

	if e.Info != nil {
		notNilCommands++
	}

	if e.Poll != nil {
		notNilCommands++
	}

	if e.Transfer != nil {
		notNilCommands++
	}

	if e.Create != nil {
		notNilCommands++
	}

	if e.Delete != nil {
		notNilCommands++
	}

	if e.Renew != nil {
		notNilCommands++
	}

	if notNilCommands != 1 {
		return errors.New("exactly one command must be present")
	}

	return e.validate()
}

func (e *EPP) validate() error {
	if e.Hello != nil {
		return e.Hello.Validate()
	}

	if e.Login != nil {
		return e.Login.Validate()
	}

	if e.Logout != nil {
		return e.Logout.Validate()
	}

	if e.Check != nil {
		return e.Check.Validate()
	}

	if e.Info != nil {
		return e.Info.Validate()
	}

	if e.Poll != nil {
		return e.Poll.Validate()
	}

	if e.Transfer != nil {
		return e.Transfer.Validate()
	}

	if e.Create != nil {
		return e.Create.Validate()
	}

	if e.Delete != nil {
		return e.Delete.Validate()
	}

	if e.Renew != nil {
		return e.Renew.Validate()
	}

	return nil
}

func (e *EPP) Command() (command2.Commander, error) {
	if e.Hello != nil {
		return e.Hello, nil
	}

	if e.Login != nil {
		return e.Login, nil
	}

	if e.Logout != nil {
		return e.Logout, nil
	}

	if e.Check != nil {
		return e.Check, nil
	}

	if e.Info != nil {
		return e.Info, nil
	}

	if e.Poll != nil {
		return e.Poll, nil
	}

	if e.Transfer != nil {
		return e.Transfer, nil
	}

	if e.Create != nil {
		return e.Create, nil
	}

	if e.Delete != nil {
		return e.Delete, nil
	}

	if e.Renew != nil {
		return e.Renew, nil
	}

	return nil, errors.New("unknown command")
}
