package session

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id"`
	UserEmail    string    `json:"user_email"`
	RefreshToken string    `json:"refresh_token"`
	IsRevoked    bool      `json:"is_revoked"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type Repository interface {
	CreateRefreshToken(ctx context.Context, s *Session) error
	GetRefreshToken(ctx context.Context, refreshToken string) (*Session, error)
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}
