package auth

import (
	"context"

	"github.com/Milua25/go-job-application-tracker/internal/user"
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	Create(ctx context.Context, u *user.User) error
}

// type Token struct {
// 	ID               uuid.UUID `json:"id"`
// 	UserID           uuid.UUID `json:"user_id"`
// 	Token            string    `json:"token"`
// 	RefreshToken     string    `json:"refresh_token"`
// 	ExpiresAt        time.Time `json:"expires_at"`
// 	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
// 	Revoked          bool      `json:"revoked"`
// }
