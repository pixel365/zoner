package auth

import "errors"

var (
	ErrInternalError        = errors.New("internal error")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrRegistrarIsBlocked   = errors.New("registrar is blocked")
	ErrCannotChangePassword = errors.New("cannot change password")
)
