package contact

import (
	"encoding/xml"

	"github.com/pixel365/goepp/command"
	"github.com/pixel365/goepp/command/update"
)

type ContactCreateResData struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:contact-1.0 creData"`
	ID      string   `xml:"urn:ietf:params:xml:ns:contact-1.0 id"`
	CRDate  string   `xml:"urn:ietf:params:xml:ns:contact-1.0 crDate"`
}

type CheckContactID struct {
	Value     string `xml:",chardata"`
	Available uint8  `xml:"avail,attr"`
}

type SingleCheckContact struct {
	Reason *string        `xml:"reason,omitempty"`
	ID     CheckContactID `xml:"urn:ietf:params:xml:ns:contact-1.0 id"`
}

type ContactsCheckResData struct {
	XMLName  xml.Name             `xml:"urn:ietf:params:xml:ns:contact-1.0 chkData"`
	Contacts []SingleCheckContact `xml:"urn:ietf:params:xml:ns:contact-1.0 cd"`
}

type ContactInfoPostalInfo struct {
	Type string              `xml:"type,attr"`
	Name string              `xml:"urn:ietf:params:xml:ns:contact-1.0 name"`
	Org  string              `xml:"urn:ietf:params:xml:ns:contact-1.0 org,omitempty"`
	Addr ContactInfoAddrData `xml:"urn:ietf:params:xml:ns:contact-1.0 addr"`
}

type ContactInfoAddrData struct {
	City    string   `xml:"urn:ietf:params:xml:ns:contact-1.0 city"`
	Sp      string   `xml:"urn:ietf:params:xml:ns:contact-1.0 sp,omitempty"`
	Pc      string   `xml:"urn:ietf:params:xml:ns:contact-1.0 pc,omitempty"`
	Cc      string   `xml:"urn:ietf:params:xml:ns:contact-1.0 cc"`
	Streets []string `xml:"urn:ietf:params:xml:ns:contact-1.0 street,omitempty"`
}

type ContactInfoDisclose struct {
	Name  *ContactInfoDisclosePostal `xml:"urn:ietf:params:xml:ns:contact-1.0 name,omitempty"`
	Org   *ContactInfoDisclosePostal `xml:"urn:ietf:params:xml:ns:contact-1.0 org,omitempty"`
	Addr  *ContactInfoDisclosePostal `xml:"urn:ietf:params:xml:ns:contact-1.0 addr,omitempty"`
	Voice *struct{}                  `xml:"urn:ietf:params:xml:ns:contact-1.0 voice,omitempty"`
	Fax   *struct{}                  `xml:"urn:ietf:params:xml:ns:contact-1.0 fax,omitempty"`
	Email *struct{}                  `xml:"urn:ietf:params:xml:ns:contact-1.0 email,omitempty"`
	Flag  uint8                      `xml:"flag,attr"`
}

type ContactInfoDisclosePostal struct {
	Type string `xml:"type,attr,omitempty"`
}

type ContactInfoResData struct {
	Voice      *string                 `xml:"urn:ietf:params:xml:ns:contact-1.0 voice,omitempty"`
	Disclose   *ContactInfoDisclose    `xml:"urn:ietf:params:xml:ns:contact-1.0 disclose,omitempty"`
	AuthInfo   *command.AuthInfo       `xml:"urn:ietf:params:xml:ns:contact-1.0 authInfo,omitempty"`
	Fax        *string                 `xml:"urn:ietf:params:xml:ns:contact-1.0 fax,omitempty"`
	XMLName    xml.Name                `xml:"urn:ietf:params:xml:ns:contact-1.0 infData"`
	UpDate     string                  `xml:"urn:ietf:params:xml:ns:contact-1.0 upDate,omitempty"`
	Email      string                  `xml:"urn:ietf:params:xml:ns:contact-1.0 email"`
	ClID       string                  `xml:"urn:ietf:params:xml:ns:contact-1.0 clID"`
	CrID       string                  `xml:"urn:ietf:params:xml:ns:contact-1.0 crID,omitempty"`
	CrDate     string                  `xml:"urn:ietf:params:xml:ns:contact-1.0 crDate,omitempty"`
	UpID       string                  `xml:"urn:ietf:params:xml:ns:contact-1.0 upID,omitempty"`
	TrDate     string                  `xml:"urn:ietf:params:xml:ns:contact-1.0 trDate,omitempty"`
	Roid       string                  `xml:"urn:ietf:params:xml:ns:contact-1.0 roid"`
	ID         string                  `xml:"urn:ietf:params:xml:ns:contact-1.0 id"`
	PostalInfo []ContactInfoPostalInfo `xml:"urn:ietf:params:xml:ns:contact-1.0 postalInfo,omitempty"`
	Statuses   []update.ContactStatus  `xml:"urn:ietf:params:xml:ns:contact-1.0 status,omitempty"`
}

//nolint:gocognit,gocyclo,cyclop
func (d ContactInfoResData) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	start := xml.StartElement{
		Name: xml.Name{Local: "contact:infData"},
		Attr: []xml.Attr{
			{
				Name:  xml.Name{Local: "xmlns:contact"},
				Value: "urn:ietf:params:xml:ns:contact-1.0",
			},
		},
	}

	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	if err := enc.EncodeElement(d.ID, xml.StartElement{Name: xml.Name{Local: "contact:id"}}); err != nil {
		return err
	}

	if err := enc.EncodeElement(d.Roid, xml.StartElement{Name: xml.Name{Local: "contact:roid"}}); err != nil {
		return err
	}

	for _, status := range d.Statuses {
		if err := enc.EncodeElement(status, xml.StartElement{Name: xml.Name{Local: "contact:status"}}); err != nil {
			return err
		}
	}

	for i := range d.PostalInfo {
		if err := enc.Encode(d.PostalInfo[i]); err != nil {
			return err
		}
	}

	if d.Voice != nil {
		if err := enc.EncodeElement(*d.Voice, xml.StartElement{Name: xml.Name{Local: "contact:voice"}}); err != nil {
			return err
		}
	}

	if d.Fax != nil {
		if err := enc.EncodeElement(*d.Fax, xml.StartElement{Name: xml.Name{Local: "contact:fax"}}); err != nil {
			return err
		}
	}

	if err := enc.EncodeElement(d.Email, xml.StartElement{Name: xml.Name{Local: "contact:email"}}); err != nil {
		return err
	}

	if err := enc.EncodeElement(d.ClID, xml.StartElement{Name: xml.Name{Local: "contact:clID"}}); err != nil {
		return err
	}

	if d.CrID != "" {
		if err := enc.EncodeElement(d.CrID, xml.StartElement{Name: xml.Name{Local: "contact:crID"}}); err != nil {
			return err
		}
	}

	if d.CrDate != "" {
		if err := enc.EncodeElement(d.CrDate, xml.StartElement{Name: xml.Name{Local: "contact:crDate"}}); err != nil {
			return err
		}
	}

	if d.UpID != "" {
		if err := enc.EncodeElement(d.UpID, xml.StartElement{Name: xml.Name{Local: "contact:upID"}}); err != nil {
			return err
		}
	}

	if d.UpDate != "" {
		if err := enc.EncodeElement(d.UpDate, xml.StartElement{Name: xml.Name{Local: "contact:upDate"}}); err != nil {
			return err
		}
	}

	if d.TrDate != "" {
		if err := enc.EncodeElement(d.TrDate, xml.StartElement{Name: xml.Name{Local: "contact:trDate"}}); err != nil {
			return err
		}
	}

	if d.AuthInfo != nil {
		if err := encodeContactAuthInfo(enc, d.AuthInfo); err != nil {
			return err
		}
	}

	if d.Disclose != nil {
		if err := enc.Encode(d.Disclose); err != nil {
			return err
		}
	}

	if err := enc.EncodeToken(start.End()); err != nil {
		return err
	}

	return enc.Flush()
}

func (p ContactInfoPostalInfo) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	start := xml.StartElement{
		Name: xml.Name{Local: "contact:postalInfo"},
		Attr: []xml.Attr{
			{
				Name:  xml.Name{Local: "type"},
				Value: p.Type,
			},
		},
	}

	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	if err := enc.EncodeElement(p.Name, xml.StartElement{Name: xml.Name{Local: "contact:name"}}); err != nil {
		return err
	}

	if p.Org != "" {
		if err := enc.EncodeElement(p.Org, xml.StartElement{Name: xml.Name{Local: "contact:org"}}); err != nil {
			return err
		}
	}

	if err := enc.Encode(p.Addr); err != nil {
		return err
	}

	if err := enc.EncodeToken(start.End()); err != nil {
		return err
	}

	return enc.Flush()
}

func (a ContactInfoAddrData) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	start := xml.StartElement{Name: xml.Name{Local: "contact:addr"}}
	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	for _, street := range a.Streets {
		if err := enc.EncodeElement(street, xml.StartElement{Name: xml.Name{Local: "contact:street"}}); err != nil {
			return err
		}
	}

	if err := enc.EncodeElement(a.City, xml.StartElement{Name: xml.Name{Local: "contact:city"}}); err != nil {
		return err
	}

	if a.Sp != "" {
		if err := enc.EncodeElement(a.Sp, xml.StartElement{Name: xml.Name{Local: "contact:sp"}}); err != nil {
			return err
		}
	}

	if a.Pc != "" {
		if err := enc.EncodeElement(a.Pc, xml.StartElement{Name: xml.Name{Local: "contact:pc"}}); err != nil {
			return err
		}
	}

	if err := enc.EncodeElement(a.Cc, xml.StartElement{Name: xml.Name{Local: "contact:cc"}}); err != nil {
		return err
	}

	if err := enc.EncodeToken(start.End()); err != nil {
		return err
	}

	return enc.Flush()
}

func (d ContactInfoDisclose) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	start := xml.StartElement{
		Name: xml.Name{Local: "contact:disclose"},
		Attr: []xml.Attr{
			{
				Name:  xml.Name{Local: "flag"},
				Value: string('0' + rune(d.Flag)),
			},
		},
	}

	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	if d.Name != nil {
		if err := encodeDisclosePostal(enc, "contact:name", d.Name.Type); err != nil {
			return err
		}
	}

	if d.Org != nil {
		if err := encodeDisclosePostal(enc, "contact:org", d.Org.Type); err != nil {
			return err
		}
	}

	if d.Addr != nil {
		if err := encodeDisclosePostal(enc, "contact:addr", d.Addr.Type); err != nil {
			return err
		}
	}

	if d.Voice != nil {
		if err := enc.EncodeElement("", xml.StartElement{Name: xml.Name{Local: "contact:voice"}}); err != nil {
			return err
		}
	}

	if d.Fax != nil {
		if err := enc.EncodeElement("", xml.StartElement{Name: xml.Name{Local: "contact:fax"}}); err != nil {
			return err
		}
	}

	if d.Email != nil {
		if err := enc.EncodeElement("", xml.StartElement{Name: xml.Name{Local: "contact:email"}}); err != nil {
			return err
		}
	}

	if err := enc.EncodeToken(start.End()); err != nil {
		return err
	}

	return enc.Flush()
}

func encodeDisclosePostal(enc *xml.Encoder, name, typ string) error {
	start := xml.StartElement{Name: xml.Name{Local: name}}
	if typ != "" {
		start.Attr = append(start.Attr, xml.Attr{
			Name:  xml.Name{Local: "type"},
			Value: typ,
		})
	}

	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	if err := enc.EncodeToken(start.End()); err != nil {
		return err
	}

	return nil
}

func encodeContactAuthInfo(enc *xml.Encoder, authInfo *command.AuthInfo) error {
	start := xml.StartElement{Name: xml.Name{Local: "contact:authInfo"}}
	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	if authInfo.Password != "" {
		if err := enc.EncodeElement(authInfo.Password, xml.StartElement{Name: xml.Name{Local: "contact:pw"}}); err != nil {
			return err
		}
	}

	if authInfo.Null != nil {
		if err := enc.EncodeElement("", xml.StartElement{Name: xml.Name{Local: "contact:null"}}); err != nil {
			return err
		}
	}

	if err := enc.EncodeToken(start.End()); err != nil {
		return err
	}

	return nil
}

func (d ContactsCheckResData) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	start := xml.StartElement{
		Name: xml.Name{Local: "contact:chkData"},
		Attr: []xml.Attr{
			{
				Name:  xml.Name{Local: "xmlns:contact"},
				Value: "urn:ietf:params:xml:ns:contact-1.0",
			},
		},
	}

	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	for _, contact := range d.Contacts {
		if err := enc.Encode(contact); err != nil {
			return err
		}
	}

	if err := enc.EncodeToken(start.End()); err != nil {
		return err
	}

	return enc.Flush()
}

func (c SingleCheckContact) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	start := xml.StartElement{Name: xml.Name{Local: "contact:cd"}}
	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	if err := enc.Encode(c.ID); err != nil {
		return err
	}

	if c.Reason != nil {
		if err := enc.EncodeElement(*c.Reason, xml.StartElement{Name: xml.Name{Local: "contact:reason"}}); err != nil {
			return err
		}
	}

	if err := enc.EncodeToken(start.End()); err != nil {
		return err
	}

	return enc.Flush()
}

func (i CheckContactID) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	start := xml.StartElement{
		Name: xml.Name{Local: "contact:id"},
		Attr: []xml.Attr{
			{
				Name:  xml.Name{Local: "avail"},
				Value: string('0' + rune(i.Available)),
			},
		},
	}

	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	if err := enc.EncodeToken(xml.CharData(i.Value)); err != nil {
		return err
	}

	if err := enc.EncodeToken(start.End()); err != nil {
		return err
	}

	return enc.Flush()
}
