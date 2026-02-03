package command

import (
	"encoding/xml"
	"errors"

	"github.com/pixel365/zoner/epp/server/internal/command/hello"
	"github.com/pixel365/zoner/epp/server/internal/command/login"
)

type EPP struct {
	// Hello https://datatracker.ietf.org/doc/html/rfc5730#section-2.3
	Hello *hello.Hello `xml:"urn:ietf:params:xml:ns:epp-1.0 hello"`

	// Login https://datatracker.ietf.org/doc/html/rfc5730#section-2.9.1.1
	Login *login.Login `xml:"urn:ietf:params:xml:ns:epp-1.0 command>login"`

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

	return nil
}

func (e *EPP) Command() (Command, error) {
	if e.Hello != nil {
		return e.Hello, nil
	}

	if e.Login != nil {
		return e.Login, nil
	}

	return nil, errors.New("unknown command")
}
