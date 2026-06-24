package token

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"           json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index"       json:"user_id"`
	RefreshToken string    `gorm:"uniqueIndex;not null;size:512"  json:"refresh_token"`
	IsRevoked    bool      `gorm:"not null;default:false"         json:"is_revoked"`
	CreatedAt    time.Time `                                      json:"created_at"`
	ExpiresAt    time.Time `gorm:"not null"                       json:"expires_at"`
}

type Repository interface {
	CreateSession(ctx context.Context, s *Session) error
	GetSessionByID(ctx context.Context, id string) (*Session, error)
	RevokeSessionByToken(ctx context.Context, token string) error
	DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteSessionByID(ctx context.Context, id string) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*Session, error)
}
