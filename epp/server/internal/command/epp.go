package command

import (
	"encoding/xml"
	"errors"
)

type EPP struct {
	Hello   *Hello   `xml:"urn:ietf:params:xml:ns:epp-1.0 hello"`
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
}

func (e EPP) Command() (Command, error) {
	if e.Hello != nil {
		return e.Hello, nil
	}

	return nil, errors.New("unknown command")
}
