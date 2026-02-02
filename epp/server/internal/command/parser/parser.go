package parser

import (
	"encoding/xml"
	"fmt"

	"github.com/pixel365/zoner/epp/server/internal/command"
)

type Parser interface {
	Parse([]byte) (command.Command, error)
}

type CommandParser struct{}

func (c *CommandParser) Parse(payload []byte) (command.Command, error) {
	var message command.EPP
	if err := xml.Unmarshal(payload, &message); err != nil {
		return nil, fmt.Errorf("unmarshal xml payload: %w", err)
	}

	return message.Command()
}
