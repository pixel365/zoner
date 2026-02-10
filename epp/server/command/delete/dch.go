package delete

func (d *Delete) SetDomain(dm *DomainDelete) {
	d.Domain = dm
}

func (d *Delete) SetContact(c *ContactDelete) {
	d.Contact = c
}

func (d *Delete) SetHost(h *HostDelete) {
	d.Host = h
}

func (d *Delete) GetDomain() *DomainDelete {
	return d.Domain
}

func (d *Delete) GetContact() *ContactDelete {
	return d.Contact
}

func (d *Delete) GetHost() *HostDelete {
	return d.Host
}
