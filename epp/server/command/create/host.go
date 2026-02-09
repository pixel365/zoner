package create

import "errors"

type HostAddrType string

const (
	HostAddrTypeIPv4 HostAddrType = "v4"
	HostAddrTypeIPv6 HostAddrType = "v6"
)

type Host struct {
	Name  string
	Addrs []hostAddrXML
}

func (h *Host) Validate() error {
	if h.Name == "" {
		return errors.New("host:name is required")
	}

	for _, a := range h.Addrs {
		if a.IP != "" && a.IP != HostAddrTypeIPv4 && a.IP != HostAddrTypeIPv6 {
			return errors.New("host:addr ip must be v4 or v6")
		}
		if a.Value == "" {
			return errors.New("host:addr value is empty")
		}
	}

	return nil
}
