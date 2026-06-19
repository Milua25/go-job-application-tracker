package user

// 👈 Defines the INTERFACE

import "context"

// Repository defines what the user package needs from a database.
// Notice we accept and return PURE domain entities (User), not DB-specific models.
type Repository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, u *User) error
}
