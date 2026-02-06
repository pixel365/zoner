package transfer

import (
	"errors"
	"fmt"

	"github.com/pixel365/zoner/epp/server/command/command"
)

type Transfer struct {
	Domain  *Domain  `xml:"-"`
	Contact *Contact `xml:"-"`
	Op      string   `xml:"op,attr"`
}

func (t *Transfer) Name() command.CommandName {
	return command.Transfer
}

func (t *Transfer) NeedAuth() bool {
	return true
}

func (t *Transfer) Validate() error {
	if t.Op == "" {
		return errors.New("transfer operation is empty")
	}

	switch t.Op {
	case "request", "query", "cancel", "approve", "reject":
	default:
		return fmt.Errorf("transfer operation %s is not supported", t.Op)
	}

	var notNil uint8

	if t.Domain != nil {
		notNil++
	}

	if t.Contact != nil {
		notNil++
	}

	if notNil != 1 {
		return errors.New("transfer must contain exactly one of domain or contact")
	}

	if t.Domain != nil {
		return t.Domain.Validate(t.Op)
	}

	if t.Contact != nil {
		return t.Contact.Validate(t.Op)
	}

	return nil
}
