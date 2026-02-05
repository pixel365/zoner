package login

import (
	"errors"

	"github.com/pixel365/zoner/epp/server/command/command"
)

type Login struct {
	Options     Options `xml:"options"`
	ClientID    string  `xml:"clID"`
	Password    string  `xml:"pw"`
	NewPassword string  `xml:"newPW,omitempty"`
	Svc         Svc     `xml:"svcs"`
}

func (l Login) Name() command.CommandName {
	return command.Login
}

func (l Login) NeedAuth() bool {
	return false
}

func (l Login) Validate() error {
	if l.ClientID == "" {
		return errors.New("client id is empty")
	}

	if l.Password == "" {
		return errors.New("password is empty")
	}

	if l.NewPassword != "" {
		if l.NewPassword == l.Password {
			return errors.New("new password must be different from old one")
		}
	}

	if err := l.Options.Validate(); err != nil {
		return err
	}

	if err := l.Svc.Validate(); err != nil {
		return err
	}

	return nil
}
