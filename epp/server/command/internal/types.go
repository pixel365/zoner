package internal

type PeriodUnit string

const (
	PeriodUnitYear  PeriodUnit = "y"
	PeriodUnitMonth PeriodUnit = "m"
)

type Period struct {
	Unit  PeriodUnit `xml:"unit,attr,omitempty"`
	Value int        `xml:",chardata"`
}
