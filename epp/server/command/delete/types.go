package delete

import "github.com/pixel365/zoner/epp/server/command/internal"

type DomainDelete struct {
	internal.DomainRef
}

type ContactDelete struct {
	internal.ContactRef
}

type HostDelete struct {
	internal.HostRef
}
