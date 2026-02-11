package update

import (
	"errors"

	"github.com/pixel365/zoner/epp/server/command/command"
)

type Update struct {
	Domain *Domain
}

func (u *Update) Name() command.CommandName {
	return command.Update
}

func (u *Update) NeedAuth() bool {
	return true
}

func (u *Update) Validate() error {
	if u.Domain == nil {
		return errors.New("exactly one update object must be present")
	}

	if u.Domain.Add == nil && u.Domain.Remove == nil && u.Domain.Change == nil {
		return errors.New("at least one update object must be present")
	}

	return u.Domain.Validate()
}
