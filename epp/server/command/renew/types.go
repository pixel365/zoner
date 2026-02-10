package renew

import "github.com/pixel365/zoner/epp/server/command/internal"

type Domain struct {
	Period *internal.Period `xml:"period,omitempty"`
	internal.DomainRef
	CurrentExpirationDate string `xml:"curExpDate"`
}
