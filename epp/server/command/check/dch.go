package check

func (c *Check) SetDomain(d *DomainCheck) {
	c.Domain = d
}

func (c *Check) SetContact(ct *ContactCheck) {
	c.Contact = ct
}

func (c *Check) SetHost(h *HostCheck) {
	c.Host = h
}

func (c *Check) GetDomain() *DomainCheck {
	return c.Domain
}

func (c *Check) GetContact() *ContactCheck {
	return c.Contact
}

func (c *Check) GetHost() *HostCheck {
	return c.Host
}
