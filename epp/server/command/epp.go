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
	"github.com/pixel365/zoner/epp/server/command/update"
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

	// Update https://datatracker.ietf.org/doc/html/rfc5730#section-2.9.3.5
	Update *update.Update `xml:"urn:ietf:params:xml:ns:epp-1.0 command>update"`

	// XMLName https://datatracker.ietf.org/doc/html/rfc5730#section-2.2
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`

	TransactionID string `xml:"urn:ietf:params:xml:ns:epp-1.0 command>clTRID,omitempty"` //TODO: add validator
}

func (e *EPP) Validate() error { //TODO: change to error response
	var notNilCommands uint8

	switch {
	case e.Hello != nil:
		notNilCommands++
	case e.Login != nil:
		notNilCommands++
	case e.Logout != nil:
		notNilCommands++
	case e.Check != nil:
		notNilCommands++
	case e.Info != nil:
		notNilCommands++
	case e.Poll != nil:
		notNilCommands++
	case e.Transfer != nil:
		notNilCommands++
	case e.Create != nil:
		notNilCommands++
	case e.Delete != nil:
		notNilCommands++
	case e.Renew != nil:
		notNilCommands++
	case e.Update != nil:
		notNilCommands++
	}

	if notNilCommands != 1 {
		return errors.New("exactly one command must be present")
	}

	return e.validate()
}

func (e *EPP) validate() error {
	switch {
	case e.Hello != nil:
		return e.Hello.Validate()
	case e.Login != nil:
		return e.Login.Validate()
	case e.Logout != nil:
		return e.Logout.Validate()
	case e.Check != nil:
		return e.Check.Validate()
	case e.Info != nil:
		return e.Info.Validate()
	case e.Poll != nil:
		return e.Poll.Validate()
	case e.Transfer != nil:
		return e.Transfer.Validate()
	case e.Create != nil:
		return e.Create.Validate()
	case e.Delete != nil:
		return e.Delete.Validate()
	case e.Renew != nil:
		return e.Renew.Validate()
	case e.Update != nil:
		return e.Update.Validate()
	default:
		return nil
	}
}

func (e *EPP) Command() (command2.Commander, error) {
	switch {
	case e.Hello != nil:
		return e.Hello, nil
	case e.Login != nil:
		return e.Login, nil
	case e.Logout != nil:
		return e.Logout, nil
	case e.Check != nil:
		return e.Check, nil
	case e.Info != nil:
		return e.Info, nil
	case e.Poll != nil:
		return e.Poll, nil
	case e.Transfer != nil:
		return e.Transfer, nil
	case e.Create != nil:
		return e.Create, nil
	case e.Delete != nil:
		return e.Delete, nil
	case e.Renew != nil:
		return e.Renew, nil
	case e.Update != nil:
		return e.Update, nil
	default:
		return nil, errors.New("unknown command")
	}
}
