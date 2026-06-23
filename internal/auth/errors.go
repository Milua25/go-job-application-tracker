package auth

import "errors"

var (
	// ErrInvalidCredentials is returned when the provided credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserNotFound is returned when a user is not found in the database.
	ErrUserNotFound = errors.New("user not found")

	// ErrTokenCreationFailed is returned when the creation of a token fails.
	ErrTokenCreationFailed = errors.New("failed to create token")

	// ErrSessionCreationFailed is returned when the creation of a session fails.
	ErrSessionCreationFailed = errors.New("failed to create session")

	// ErrUserAlreadyExists is returned when a user with the same email already exists.
	ErrUserAlreadyExists = errors.New("user with this email already exists")

	// ErrNotFound is returned when a requested resource is not found.
	ErrNotFound = errors.New("user not found")

	ErrInvalidRefreshToken = errors.New("invalid refresh token")

	ErrFailedToDeleteSession = errors.New("failed to delete session")
)
