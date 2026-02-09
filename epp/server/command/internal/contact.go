package internal

import (
	"encoding/xml"
	"errors"
)

type Contact struct {
	AuthInfo *AuthInfo `xml:"authInfo,omitempty"`
	XMLName  xml.Name  `xml:"urn:ietf:params:xml:ns:contact-1.0 info"`
	Id       string    `xml:"id"`
}

func (c Contact) Validate() error {
	if c.Id == "" {
		return errors.New("contact:id is required")
	}

	if c.AuthInfo != nil && c.AuthInfo.Password == "" {
		return errors.New("contact:authInfo:pw is required")
	}

	return nil
}
