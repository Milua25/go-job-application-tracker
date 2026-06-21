package auth

import (
	"context"

	"github.com/Milua25/go-job-application-tracker/internal/user"
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	Create(ctx context.Context, u *user.User) error
}
