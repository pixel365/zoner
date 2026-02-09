package internal

import (
	"encoding/xml"
	"errors"
)

type Domain struct {
	AuthInfo *AuthInfo `xml:"authInfo,omitempty"`
	XMLName  xml.Name  `xml:"urn:ietf:params:xml:ns:domain-1.0 info"`
	Name     string    `xml:"name"`
}

func (d Domain) Validate() error {
	if d.Name == "" {
		return errors.New("domain:name is required")
	}

	if d.AuthInfo != nil && d.AuthInfo.Password == "" {
		return errors.New("domain:authInfo:pw is required")
	}

	return nil
}
