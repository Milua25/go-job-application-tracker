package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id string) error
	FindAllWithSessions(ctx context.Context) ([]*User, error)
}

type SessionRevoker interface {
	DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error
}
