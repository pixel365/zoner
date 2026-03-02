package repository

import "context"

type Authenticator interface {
	Login(context.Context, string, string) error
	Logout() error
}
