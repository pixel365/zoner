package check

import (
	"encoding/xml"
	"errors"
)

type Check struct {
	Domains []Domain
}

type Domain struct {
	Name string `xml:",chardata"`
}

func (c *Check) Name() string {
	return "check"
}

func (c *Check) Validate() error {
	if len(c.Domains) == 0 {
		return errors.New("objects is empty")
	}

	return nil
}

func (c *Check) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var tmp struct {
		Obj struct {
			Domains []Domain `xml:"name"`
		} `xml:"urn:ietf:params:xml:ns:obj check"`
	}

	if err := d.DecodeElement(&tmp, &start); err != nil {
		return err
	}

	domainNames := make([]Domain, 0, len(tmp.Obj.Domains))

	for _, n := range tmp.Obj.Domains {
		if n.Name != "" {
			//TODO: validate name
			domainNames = append(domainNames, n)
		}
	}

	c.Domains = domainNames

	return c.Validate()
}
