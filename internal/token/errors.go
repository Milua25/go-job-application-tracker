package token

import "errors"

// Define custom errors for token validation
var (
	ErrTokenExpired           = errors.New("token has expired")
	ErrInvalidToken           = errors.New("invalid token")
	ErrTokenEmpty             = errors.New("token is empty")
	ErrServiceNotInitialized  = errors.New("auth service is not properly initialized")
	ErrValidation             = errors.New("error validating token")
	ErrTokenMalformed         = errors.New("malformed token")
	ErrIncorrectSigningMethod = errors.New("incorrect signing method")
	ErrNotFound               = errors.New("session not found")
)
