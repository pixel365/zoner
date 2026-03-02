package errors

import "errors"

var (
	ErrInternalError      = errors.New("internal error")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
