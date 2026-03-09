package contact

import "encoding/xml"

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
