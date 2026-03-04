package model

type DomainCreateInput struct {
	Name         string
	Punycode     string
	AuthInfoHash string
	PeriodUnit   string
	RegistrarID  int64
	Period       int
}
