package check

type DomainCheck struct {
	Names []string `xml:"name"`
}

type ContactCheck struct {
	IDs []string `xml:"id"`
}

type HostCheck struct {
	Names []string `xml:"name"`
}
