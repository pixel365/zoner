package transfer

import (
	"encoding/xml"
	"errors"

	"github.com/pixel365/zoner/epp/server/command/internal"
)

type Contact struct {
	AuthInfo *internal.AuthInfo `xml:"authInfo,omitempty"`
	ID       string             `xml:"id"`
}

func (c *Contact) Validate(_ string) error {
	if c.ID == "" {
		return errors.New("contact:id is required")
	}

	//TODO: if op == request: check AuthInfo? read enable auth require from config

	return nil
}

type contactTransferXML struct {
	AuthInfo *internal.AuthInfo `xml:"authInfo,omitempty"`
	XMLName  xml.Name           `xml:"urn:ietf:params:xml:ns:contact-1.0 transfer"`
	ID       string             `xml:"id"`
}
