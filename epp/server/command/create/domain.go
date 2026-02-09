package create

import (
	"encoding/xml"
	"errors"

	"github.com/pixel365/zoner/epp/server/command/internal"
)

type PeriodUnit string
type DomainContactType string

const (
	PeriodUnitYear  PeriodUnit = "y"
	PeriodUnitMonth PeriodUnit = "m"

	DomainContactTypeAdmin   DomainContactType = "admin"
	DomainContactTypeTech    DomainContactType = "tech"
	DomainContactTypeBilling DomainContactType = "billing"
)

type Domain struct {
	Period     *Period
	NS         *DomainNS
	Name       string
	Registrant string
	AuthInfo   internal.AuthInfo
	Contacts   []DomainContact
}

type Period struct {
	Unit  PeriodUnit
	Value int
}

type DomainNS struct {
	HostObj  []string
	HostAttr []HostAttr
}

type HostAttr struct {
	Name  string
	Addrs []domainHostAddrXML
}

type DomainContact struct {
	Type DomainContactType
	ID   string
}

type domainCreateXML struct {
	Period     *domainPeriodXML   `xml:"period,omitempty"`
	NS         *domainNsXML       `xml:"ns,omitempty"`
	XMLName    xml.Name           `xml:"urn:ietf:params:xml:ns:domain-1.0 create"`
	Name       string             `xml:"name"`
	Registrant string             `xml:"registrant,omitempty"`
	AuthInfo   domainAuthInfoXML  `xml:"authInfo"`
	Contacts   []domainContactXML `xml:"contact,omitempty"`
}

type domainPeriodXML struct {
	Unit  PeriodUnit `xml:"unit,attr,omitempty"`
	Value int        `xml:",chardata"`
}

type domainNsXML struct {
	HostObj  []string            `xml:"hostObj,omitempty"`
	HostAttr []domainHostAttrXML `xml:"hostAttr,omitempty"`
}

type domainHostAttrXML struct {
	Name  string              `xml:"hostName"`
	Addrs []domainHostAddrXML `xml:"hostAddr,omitempty"`
}

type domainHostAddrXML struct {
	IP    string `xml:"ip,attr,omitempty"`
	Value string `xml:",chardata"`
}

type domainContactXML struct {
	Type  DomainContactType `xml:"type,attr"`
	Value string            `xml:",chardata"`
}

type domainAuthInfoXML struct {
	PW string `xml:"pw"`
}

func (d *Domain) Validate() error {
	if d.Name == "" {
		return errors.New("domain:name is required")
	}

	if d.AuthInfo.Password == "" {
		return errors.New("domain:authInfo/domain:pw is required")
	}

	if d.Period != nil {
		if d.Period.Value <= 0 {
			return errors.New("domain:period must be > 0")
		}
		if d.Period.Unit != "" && d.Period.Unit != PeriodUnitYear &&
			d.Period.Unit != PeriodUnitMonth {
			return errors.New("domain:period unit must be y or m")
		}
	}

	for _, c := range d.Contacts {
		if c.Type == "" {
			return errors.New("domain:contact type is empty")
		}
		if c.ID == "" {
			return errors.New("domain:contact id is empty")
		}

		switch c.Type {
		case DomainContactTypeAdmin, DomainContactTypeTech, DomainContactTypeBilling:
		default:
			return errors.New("domain:contact type must be admin, tech or billing")
		}
	}

	if err := d.validateNS(); err != nil {
		return err
	}

	return nil
}

func (d *Domain) validateNS() error {
	if d.NS == nil {
		return nil
	}

	if len(d.NS.HostObj) == 0 && len(d.NS.HostAttr) == 0 {
		return errors.New("domain:ns must contain hostObj or hostAttr")
	}

	for _, h := range d.NS.HostObj {
		if h == "" {
			return errors.New("domain:ns/hostObj is empty")
		}
	}

	for _, ha := range d.NS.HostAttr {
		if ha.Name == "" {
			return errors.New("domain:ns/hostAttr/hostName is empty")
		}
		for _, a := range ha.Addrs {
			if a.IP != "" && a.IP != "v4" && a.IP != "v6" {
				return errors.New("domain:ns/hostAttr/hostAddr ip must be v4 or v6")
			}
			if a.Value == "" {
				return errors.New("domain:ns/hostAttr/hostAddr value is empty")
			}
		}
	}

	return nil
}
