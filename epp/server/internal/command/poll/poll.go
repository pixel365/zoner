package poll

import "fmt"

type Poll struct {
	Op string `xml:"op,attr"`
}

func (p *Poll) Name() string {
	return "poll"
}

func (p *Poll) Validate() error {
	switch p.Op {
	case "req", "ack":
		return nil
	default:
		return fmt.Errorf("invalid poll operation: %s", p.Op)
	}
}
