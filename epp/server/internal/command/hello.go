package command

import "encoding/xml"

type Hello struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:epp-1.0 hello"`
}

func (h Hello) Name() Name {
	return HelloCommand
}

func (h Hello) ClTRID() string {
	return ""
}

func (h Hello) AsBytes() []byte {
	return nil
}
