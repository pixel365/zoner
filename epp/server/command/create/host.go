package create

import "errors"

type Host struct {
	Name  string
	Addrs []hostAddrXML
}

type HostAddr struct {
	IP    string
	Value string
}

func (h *Host) Validate() error {
	if h.Name == "" {
		return errors.New("host:name is required")
	}

	for _, a := range h.Addrs {
		if a.IP != "" && a.IP != "v4" && a.IP != "v6" {
			return errors.New("host:addr ip must be v4 or v6")
		}
		if a.Value == "" {
			return errors.New("host:addr value is empty")
		}
	}

	return nil
}
