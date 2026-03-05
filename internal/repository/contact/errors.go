package contact

import "errors"

var (
	ErrAlreadyExists = errors.New("contact already exists")
	ErrValidation    = errors.New("validation error")
	ErrInternal      = errors.New("internal error")
)
