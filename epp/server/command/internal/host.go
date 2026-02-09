package internal

import (
	"encoding/xml"
	"errors"
)

type Host struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:host-1.0 info"`
	Name    string   `xml:"name"`
}

func (h Host) Validate() error {
	if h.Name == "" {
		return errors.New("host:name is required")
	}

	return nil
}
