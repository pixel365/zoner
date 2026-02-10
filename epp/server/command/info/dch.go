package info

func (i *Info) SetDomain(d *DomainInfo) {
	i.Domain = d
}

func (i *Info) SetContact(c *ContactInfo) {
	i.Contact = c
}

func (i *Info) SetHost(h *HostInfo) {
	i.Host = h
}

func (i *Info) GetDomain() *DomainInfo {
	return i.Domain
}

func (i *Info) GetContact() *ContactInfo {
	return i.Contact
}

func (i *Info) GetHost() *HostInfo {
	return i.Host
}
