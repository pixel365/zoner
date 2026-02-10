package transfer

import (
	"encoding/xml"
	"errors"

	"github.com/pixel365/zoner/epp/server/command/internal"
)

type DomainTransfer struct {
	AuthInfo *internal.AuthInfo `xml:"authInfo,omitempty"`
	internal.DomainRef
}

type ContactTransfer struct {
	AuthInfo *internal.AuthInfo `xml:"authInfo,omitempty"`
	internal.ContactRef
}

type domainTransferXML struct {
	AuthInfo *internal.AuthInfo `xml:"authInfo,omitempty"`
	XMLName  xml.Name           `xml:"urn:ietf:params:xml:ns:domain-1.0 transfer"`
	Name     string             `xml:"name"`
}

type contactTransferXML struct {
	AuthInfo *internal.AuthInfo `xml:"authInfo,omitempty"`
	XMLName  xml.Name           `xml:"urn:ietf:params:xml:ns:contact-1.0 transfer"`
	ID       string             `xml:"id"`
}

func (d *DomainTransfer) Validate(_ string) error {
	if d.AuthInfo != nil {
		if d.AuthInfo.Password == "" {
			return errors.New("domain:authInfo/domain:pw is required if authInfo is present")
		}
	}

	return d.DomainRef.Validate()
}

func (c *ContactTransfer) Validate(_ string) error {
	if c.AuthInfo != nil {
		if c.AuthInfo.Password == "" {
			return errors.New("contact:authInfo/contact:pw is required if authInfo is present")
		}
	}

	return c.ContactRef.Validate()
}
