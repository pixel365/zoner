package contact

import "encoding/xml"

type ContactCreateResData struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:contact-1.0 creData"`
	ID      string   `xml:"urn:ietf:params:xml:ns:contact-1.0 id"`
	CRDate  string   `xml:"urn:ietf:params:xml:ns:contact-1.0 crDate"`
}
