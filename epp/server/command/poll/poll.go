package poll

import (
	"fmt"

	"github.com/pixel365/zoner/epp/server/command/command"
)

type Poll struct {
	Op string `xml:"op,attr"`
}

func (p *Poll) Name() command.CommandName {
	return command.Poll
}

func (p *Poll) NeedAuth() bool {
	return true
}

func (p *Poll) Validate() error {
	switch p.Op {
	case "req", "ack":
		return nil
	default:
		return fmt.Errorf("invalid poll operation: %s", p.Op)
	}
}
