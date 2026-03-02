package repository

import "context"

type Authenticator interface {
	Login(context.Context, string, string) (int64, error)
	Logout() error
}
