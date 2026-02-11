package internal

import (
	"errors"

	normalizer "github.com/pixel365/domain-normalizer"
)

type DomainRef struct {
	Name string `xml:"name"`
}
type ContactRef struct {
	ID string `xml:"id"`
}
type HostRef struct {
	Name string `xml:"name"`
}
type AuthInfo struct {
	Null     *struct{} `xml:"null,omitempty"`
	Password string    `xml:"pw"`
}

func (d *DomainRef) Validate() error {
	if d.Name == "" {
		return errors.New("domain:name is required")
	}

	name, err := normalizer.Parse(d.Name)
	if err != nil {
		return errors.New("domain:name is invalid")
	}

	d.Name = name.Normalized

	return nil
}

func (c ContactRef) Validate() error {
	if c.ID == "" {
		return errors.New("contact:id is required")
	}

	return nil
}

func (h HostRef) Validate() error {
	if h.Name == "" {
		return errors.New("host:name is required")
	}

	return nil
}

func (a AuthInfo) Validate() error {
	if a.Password == "" && a.Null != nil {
		return errors.New("authInfo/null is not allowed when pw is present")
	}

	if a.Password != "" && a.Null != nil {
		return errors.New("authInfo/pw and authInfo/null are mutually exclusive")
	}

	if a.Password == "" {
		return errors.New("authInfo/pw is required")
	}

	return nil
}
