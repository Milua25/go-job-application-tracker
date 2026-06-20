package user

// 👈 Defines the INTERFACE

import "context"

// UserRepository defines what the user package needs from a database.
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, u *User) error
}
