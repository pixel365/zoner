package info

import "github.com/pixel365/zoner/epp/server/command/internal"

func (i *Info) SetDomain(d *internal.Domain) {
	i.Domain = d
}

func (i *Info) SetContact(c *internal.Contact) {
	i.Contact = c
}

func (i *Info) SetHost(h *internal.Host) {
	i.Host = h
}

func (i *Info) GetDomain() *internal.Domain {
	return i.Domain
}

func (i *Info) GetContact() *internal.Contact {
	return i.Contact
}

func (i *Info) GetHost() *internal.Host {
	return i.Host
}
