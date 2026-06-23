package user

import "errors"

var (
	ErrNotFound        = errors.New("user not found")
	ErrEmailInUse      = errors.New("email already in use")
)
