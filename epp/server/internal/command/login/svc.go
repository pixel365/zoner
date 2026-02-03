package login

import "errors"

type Svc struct {
	Extension *SvcExtension `xml:"svcExtension,omitempty"`
	ObjectURI []string      `xml:"objURI"`
}

func (s Svc) Validate() error {
	if len(s.ObjectURI) == 0 {
		return errors.New("object uri is empty")
	}

	if s.Extension != nil {
		if len(s.Extension.ExtensionURI) == 0 {
			return errors.New("extension uri is empty")
		}
	}

	return nil
}

type SvcExtension struct {
	ExtensionURI []string `xml:"extURI"`
}
