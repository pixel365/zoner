package contact

import "errors"

var (
	ErrAlreadyExists = errors.New("contact already exists")
	ErrNotFound      = errors.New("contact not found")
	ErrValidation    = errors.New("validation error")
	ErrInternal      = errors.New("internal error")
)
