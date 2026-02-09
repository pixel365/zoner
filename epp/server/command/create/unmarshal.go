package create

import (
	"encoding/xml"
	"errors"

	"github.com/pixel365/zoner/epp/server/command/internal"
)

func (c *Create) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	c.Domain, c.Contact, c.Host = nil, nil, nil

	var seen bool
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}

		done, err := c.handleToken(d, tok, &start, &seen)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}
}

type hostCreateXML struct {
	XMLName xml.Name      `xml:"urn:ietf:params:xml:ns:host-1.0 create"`
	Name    string        `xml:"name"`
	Addrs   []hostAddrXML `xml:"addr,omitempty"`
}

type hostAddrXML struct {
	IP    string `xml:"ip,attr,omitempty"`
	Value string `xml:",chardata"`
}

func mapHostCreate(x hostCreateXML) *Host {
	h := &Host{
		Name: x.Name,
	}

	if len(x.Addrs) > 0 {
		h.Addrs = make([]hostAddrXML, 0, len(x.Addrs))
		for _, a := range x.Addrs {
			h.Addrs = append(h.Addrs, hostAddrXML{
				IP:    a.IP,
				Value: a.Value,
			})
		}
	}
	return h
}

func mapDomainCreate(x domainCreateXML) *Domain {
	d := &Domain{
		Name:       x.Name,
		Registrant: x.Registrant,
		AuthInfo:   internal.AuthInfo{Password: x.AuthInfo.PW},
	}

	if x.Period != nil {
		d.Period = &Period{
			Unit:  x.Period.Unit,
			Value: x.Period.Value,
		}
	}

	if x.NS != nil {
		ns := &DomainNS{}
		if len(x.NS.HostObj) > 0 {
			ns.HostObj = append(ns.HostObj, x.NS.HostObj...)
		}

		if len(x.NS.HostAttr) > 0 {
			ns.HostAttr = make([]HostAttr, 0, len(x.NS.HostAttr))
			for _, ha := range x.NS.HostAttr {
				h := HostAttr{Name: ha.Name}
				if len(ha.Addrs) > 0 {
					h.Addrs = make([]domainHostAddrXML, 0, len(ha.Addrs))
					for _, a := range ha.Addrs {
						h.Addrs = append(h.Addrs, domainHostAddrXML{
							IP:    a.IP,
							Value: a.Value,
						})
					}
				}
				ns.HostAttr = append(ns.HostAttr, h)
			}
		}

		d.NS = ns
	}

	if len(x.Contacts) == 0 {
		return d
	}

	d.Contacts = make([]DomainContact, 0, len(x.Contacts))
	for _, c := range x.Contacts {
		d.Contacts = append(d.Contacts, DomainContact{
			Type: c.Type,
			ID:   c.Value,
		})
	}

	return d
}

func mapContactCreate(x contactCreateXML) *Contact {
	c := &Contact{
		ID:    x.ID,
		Email: x.Email,
		AuthInfo: internal.AuthInfo{
			Password: x.AuthInfo.Password,
		},
	}

	if x.Voice != nil {
		c.Voice = &ContactPhone{X: x.Voice.X, Num: x.Voice.Num}
	}

	if x.Fax != nil {
		c.Fax = &ContactPhone{X: x.Fax.X, Num: x.Fax.Num}
	}

	if len(x.PostalInfo) > 0 {
		c.PostalInfo = make([]PostalInfo, 0, len(x.PostalInfo))
		for i := range x.PostalInfo {
			c.PostalInfo = append(c.PostalInfo, PostalInfo{
				Type: PostalInfoType(x.PostalInfo[i].Type),
				Name: x.PostalInfo[i].Name,
				Org:  x.PostalInfo[i].Org,
				Addr: ContactAddr{
					Street: append([]string(nil), x.PostalInfo[i].Addr.Street...),
					City:   x.PostalInfo[i].Addr.City,
					Sp:     x.PostalInfo[i].Addr.Sp,
					Pc:     x.PostalInfo[i].Addr.Pc,
					Cc:     x.PostalInfo[i].Addr.Cc,
				},
			})
		}
	}

	if x.Disclose != nil {
		d := &Disclose{Flag: x.Disclose.Flag}

		if x.Disclose.Name != nil {
			d.Items = append(d.Items, DiscloseItem{Name: "name", Type: x.Disclose.Name.Type})
		}

		if x.Disclose.Org != nil {
			d.Items = append(d.Items, DiscloseItem{Name: "org", Type: x.Disclose.Org.Type})
		}

		if x.Disclose.Addr != nil {
			d.Items = append(d.Items, DiscloseItem{Name: "addr", Type: x.Disclose.Addr.Type})
		}

		if x.Disclose.Voice != nil {
			d.Items = append(d.Items, DiscloseItem{Name: "voice"})
		}

		if x.Disclose.Fax != nil {
			d.Items = append(d.Items, DiscloseItem{Name: "fax"})
		}

		if x.Disclose.Email != nil {
			d.Items = append(d.Items, DiscloseItem{Name: "email"})
		}

		c.Disclose = d
	}

	return c
}

func (c *Create) handleToken(
	d *xml.Decoder,
	tok xml.Token,
	start *xml.StartElement,
	seen *bool,
) (bool, error) {
	switch t := tok.(type) {
	case xml.StartElement:
		if t.Name.Local != "create" {
			return false, errors.New("unexpected element inside <create>: " + t.Name.Local)
		}
		if *seen {
			return false, errors.New("exactly one create object must be present")
		}

		if err := c.decodeObjectCreate(d, &t); err != nil {
			return false, err
		}
		*seen = true
		return false, nil

	case xml.EndElement:
		if t.Name == start.Name {
			if !*seen {
				return false, errors.New("exactly one create object must be present")
			}
			return true, nil
		}
	}
	return false, nil
}

func (c *Create) decodeObjectCreate(d *xml.Decoder, t *xml.StartElement) error {
	switch t.Name.Space {
	case internal.NsDomain:
		var x domainCreateXML
		if err := d.DecodeElement(&x, t); err != nil {
			return err
		}
		c.Domain = mapDomainCreate(x)
		return nil

	case internal.NsContact:
		var x contactCreateXML
		if err := d.DecodeElement(&x, t); err != nil {
			return err
		}
		c.Contact = mapContactCreate(x)
		return nil

	case internal.NsHost:
		var x hostCreateXML
		if err := d.DecodeElement(&x, t); err != nil {
			return err
		}
		c.Host = mapHostCreate(x)
		return nil

	default:
		return errors.New("unsupported <create> object namespace: " + t.Name.Space)
	}
}
