package command

import (
	"encoding/xml"
	"fmt"
)

type Parser interface {
	Parse([]byte) (Command, error)
}

type CmdParser struct{}

func (c *CmdParser) Parse(payload []byte) (Command, error) {
	var message EPP
	if err := xml.Unmarshal(payload, &message); err != nil {
		return nil, fmt.Errorf("unmarshal xml payload: %w", err)
	}

	if err := message.Validate(); err != nil {
		return nil, err
	}

	return message.Command()
}
