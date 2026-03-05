package repository

import "context"

type AuthRepository interface {
	Login(context.Context, string, string, string) (int64, int64, error)
	Logout() error
}
