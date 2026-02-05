package command

import (
	"encoding/xml"
	"fmt"

	command2 "github.com/pixel365/zoner/epp/server/command/command"
)

type Parser interface {
	Parse([]byte) (command2.Commander, error)
}

type CmdParser struct{}

func (c *CmdParser) Parse(payload []byte) (command2.Commander, error) {
	var message EPP
	if err := xml.Unmarshal(payload, &message); err != nil {
		return nil, fmt.Errorf("unmarshal xml payload: %w", err)
	}

	if err := message.Validate(); err != nil {
		return nil, err
	}

	return message.Command()
}
