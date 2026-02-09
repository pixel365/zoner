package create

import (
	"encoding/xml"
	"errors"

	"github.com/pixel365/zoner/epp/server/command/internal"
)

type PostalInfoType string

const (
	Loc PostalInfoType = "loc"
	Int PostalInfoType = "int"
)

type Contact struct {
	Voice      *ContactPhone
	Fax        *ContactPhone
	Disclose   *Disclose
	ID         string
	Email      string
	AuthInfo   internal.AuthInfo
	PostalInfo []PostalInfo
}

type PostalInfo struct {
	Type PostalInfoType
	Name string
	Org  string
	Addr ContactAddr
}

type ContactAddr struct {
	City   string
	Sp     string
	Pc     string
	Cc     string
	Street []string
}

type ContactPhone struct {
	X   string
	Num string
}

type Disclose struct {
	Items []DiscloseItem
	Flag  int
}

type DiscloseItem struct {
	Name string
	Type string
}

type contactCreateXML struct {
	Voice      *phoneXML         `xml:"voice,omitempty"`
	Fax        *phoneXML         `xml:"fax,omitempty"`
	Disclose   *discloseXML      `xml:"disclose,omitempty"`
	XMLName    xml.Name          `xml:"urn:ietf:params:xml:ns:contact-1.0 create"`
	ID         string            `xml:"id"`
	Email      string            `xml:"email"`
	AuthInfo   internal.AuthInfo `xml:"authInfo"`
	PostalInfo []postalInfoXML   `xml:"postalInfo"`
}

type postalInfoXML struct {
	Type string  `xml:"type,attr"`
	Name string  `xml:"name"`
	Org  string  `xml:"org,omitempty"`
	Addr addrXML `xml:"addr"`
}

type addrXML struct {
	City   string   `xml:"city"`
	Sp     string   `xml:"sp,omitempty"`
	Pc     string   `xml:"pc,omitempty"`
	Cc     string   `xml:"cc"`
	Street []string `xml:"street,omitempty"`
}

type phoneXML struct {
	X   string `xml:"x,attr,omitempty"`
	Num string `xml:",chardata"`
}

type discloseXML struct {
	Name  *disclosePIXML `xml:"name,omitempty"`
	Org   *disclosePIXML `xml:"org,omitempty"`
	Addr  *disclosePIXML `xml:"addr,omitempty"`
	Voice *struct{}      `xml:"voice,omitempty"`
	Fax   *struct{}      `xml:"fax,omitempty"`
	Email *struct{}      `xml:"email,omitempty"`
	Flag  int            `xml:"flag,attr"`
}

type disclosePIXML struct {
	Type string `xml:"type,attr,omitempty"`
}

func (c *Contact) Validate() error {
	if c.ID == "" {
		return errors.New("contact:id is required")
	}

	if len(c.PostalInfo) == 0 {
		return errors.New("contact:postalInfo is required")
	}

	if err := c.validatePostInfo(); err != nil {
		return err
	}

	if c.Email == "" {
		return errors.New("contact:email is required")
	}

	if c.AuthInfo.Password == "" {
		return errors.New("contact:authInfo/contact:pw is required")
	}

	if c.Disclose == nil {
		return nil
	}

	if c.Disclose.Flag != 0 && c.Disclose.Flag != 1 {
		return errors.New("contact:disclose flag must be 0 or 1")
	}
	for _, it := range c.Disclose.Items {
		if it.Name == "" {
			return errors.New("contact:disclose item name is empty")
		}
		if it.Type != "" && it.Type != "loc" && it.Type != "int" {
			return errors.New("contact:disclose item type must be loc or int")
		}
	}

	return nil
}

func (c *Contact) validatePostInfo() error {
	for i := range c.PostalInfo {
		if c.PostalInfo[i].Type != Loc && c.PostalInfo[i].Type != Int {
			return errors.New("contact:postalInfo type must be loc or int")
		}

		if c.PostalInfo[i].Name == "" {
			return errors.New("contact:postalInfo/contact:name is required")
		}

		if c.PostalInfo[i].Addr.City == "" {
			return errors.New("contact:postalInfo/contact:addr/contact:city is required")
		}

		if c.PostalInfo[i].Addr.Cc == "" {
			return errors.New("contact:postalInfo/contact:addr/contact:cc is required")
		}

		if len(c.PostalInfo[i].Addr.Street) > 3 {
			return errors.New("contact:postalInfo/contact:addr/contact:street max is 3")
		}
	}

	return nil
}
