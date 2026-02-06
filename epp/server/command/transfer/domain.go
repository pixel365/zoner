package transfer

import (
	"encoding/xml"
	"errors"
)

type Domain struct {
	AuthInfo *AuthInfo `xml:"authInfo,omitempty"`
	Name     string    `xml:"name"`
}

func (d *Domain) Validate(_ string) error {
	if d.Name == "" {
		return errors.New("domain:name is required")
	}

	//TODO: if op == request: check AuthInfo? read enable auth require from config

	return nil
}

type domainTransferXML struct {
	AuthInfo *AuthInfo `xml:"authInfo,omitempty"`
	XMLName  xml.Name  `xml:"urn:ietf:params:xml:ns:domain-1.0 transfer"`
	Name     string    `xml:"name"`
}
