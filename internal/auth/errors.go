package auth

import "errors"

var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrTokenCreationFailed   = errors.New("failed to create token")
	ErrSessionCreationFailed = errors.New("failed to create session")
	ErrUserAlreadyExists     = errors.New("user with this email already exists")
	ErrNotFound              = errors.New("resource not found")
	ErrInvalidRefreshToken   = errors.New("invalid refresh token")
	ErrFailedToDeleteSession = errors.New("failed to delete sessions")
	ErrFailedToRevokeSession = errors.New("failed to revoke session")
)
