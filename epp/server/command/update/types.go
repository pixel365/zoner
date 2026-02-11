package update

import (
	"errors"

	normalizer "github.com/pixel365/domain-normalizer"

	"github.com/pixel365/zoner/epp/server/command/internal"
)

type Domain struct {
	Add    *Add    `xml:"add,omitempty"`
	Remove *Remove `xml:"rem,omitempty"`
	Change *Change `xml:"chg,omitempty"`
	internal.DomainRef
}

type Add struct {
	NS       *NS       `xml:"ns,omitempty"`
	Contacts []Contact `xml:"contact,omitempty"`
	Statuses []Status  `xml:"status,omitempty"`
}

type Remove struct {
	NS       *NS       `xml:"ns,omitempty"`
	Contacts []Contact `xml:"contact,omitempty"`
	Statuses []Status  `xml:"status,omitempty"`
}

type Change struct {
	Registrant *string            `xml:"registrant,omitempty"`
	AuthInfo   *internal.AuthInfo `xml:"authInfo,omitempty"`
}

type NS struct {
	Hosts []HostObject `xml:"hostObj,omitempty"`
}

type HostObject struct {
	Name string `xml:",chardata"`
}

type Contact struct {
	ID   string `xml:",chardata"`
	Type string `xml:"type,attr"`
}

type Status struct {
	Value string `xml:"s,attr"`
}

// Validate https://datatracker.ietf.org/doc/html/rfc5731#section-3.2.5
func (d *Domain) Validate() error {
	if d.Name == "" {
		return errors.New("domain:name is required")
	}

	name, err := normalizer.Parse(d.Name)
	if err != nil {
		return errors.New("domain:name is invalid")
	}

	d.Name = name.Normalized

	if d.Add != nil {
		if err = d.Add.Validate(); err != nil {
			return err
		}
	}

	if d.Remove != nil {
		if err = d.Remove.Validate(); err != nil {
			return err
		}
	}

	if d.Change != nil {
		if err = d.Change.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (a *Add) Validate() error {
	if a.NS == nil && len(a.Contacts) == 0 && len(a.Statuses) == 0 {
		return errors.New("add:ns, contact or status is required")
	}

	if a.NS != nil {
		if err := a.NS.Validate(); err != nil {
			return err
		}
	}

	if len(a.Contacts) > 0 {
		for _, c := range a.Contacts {
			if c.ID == "" {
				return errors.New("add:contact id is required")
			}
		}
	}

	if len(a.Statuses) > 0 {
		for _, s := range a.Statuses {
			if s.Value == "" {
				return errors.New("add:status value is required")
			}
		}
	}

	return nil
}

func (r *Remove) Validate() error {
	if r.NS == nil && len(r.Contacts) == 0 && len(r.Statuses) == 0 {
		return errors.New("rem:ns, contact or status is required")
	}

	if r.NS != nil {
		if err := r.NS.Validate(); err != nil {
			return err
		}
	}

	if len(r.Contacts) > 0 {
		for _, c := range r.Contacts {
			if c.ID == "" {
				return errors.New("rem:contact id is required")
			}
		}
	}

	if len(r.Statuses) > 0 {
		for _, s := range r.Statuses {
			if s.Value == "" {
				return errors.New("rem:status value is required")
			}
		}
	}

	return nil
}

func (c *Change) Validate() error {
	if c.Registrant == nil && c.AuthInfo == nil {
		return errors.New("change:registrant or/and authInfo is required")
	}

	if c.AuthInfo != nil {
		if err := c.AuthInfo.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (ns NS) Validate() error {
	if len(ns.Hosts) == 0 {
		return errors.New("ns:hostObj is required")
	}

	for _, h := range ns.Hosts {
		if h.Name == "" {
			return errors.New("ns:hostObj is empty")
		}
	}

	return nil
}
