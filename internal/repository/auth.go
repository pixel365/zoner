package repository

import "context"

type AuthRepository interface {
	Login(context.Context, string, string) (int64, error)
	Logout() error
}
