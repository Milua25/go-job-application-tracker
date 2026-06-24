package sqlconnect

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Milua25/go-job-application-tracker/internal/token"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionStore struct {
	db *gorm.DB
}

func (s *SessionStore) CreateSession(ctx context.Context, newUserSession *token.Session) error {
	return withTx(s.db.WithContext(ctx), func(tx *gorm.DB) error {
		return tx.Create(newUserSession).Error
	}, &sql.TxOptions{})
}

func (s *SessionStore) GetSessionByID(ctx context.Context, sessionID string) (*token.Session, error) {
	var sess token.Session
	err := s.db.WithContext(ctx).First(&sess, "id = ?", sessionID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, token.ErrNotFound
		}
		return nil, err
	}
	return &sess, nil
}

func (s *SessionStore) RevokeSessionByToken(ctx context.Context, refreshToken string) error {
	slog.Info("revoking session", "token_hash", refreshToken[:8])
	return withTx(s.db.WithContext(ctx), func(tx *gorm.DB) error {
		return tx.Model(&token.Session{}).
			Where("refresh_token = ?", refreshToken).
			Update("is_revoked", true).Error
	}, &sql.TxOptions{})
}

func (s *SessionStore) DeleteSessionByID(ctx context.Context, sessionID string) error {
	return withTx(s.db.WithContext(ctx), func(tx *gorm.DB) error {
		result := tx.Delete(&token.Session{}, "id = ?", sessionID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return token.ErrNotFound
		}
		return nil
	}, &sql.TxOptions{})
}

func (s *SessionStore) DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error {
	return withTx(s.db.WithContext(ctx), func(tx *gorm.DB) error {
		return tx.Delete(&token.Session{}, "user_id = ?", userID).Error
	}, &sql.TxOptions{})
}

func (s *SessionStore) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*token.Session, error) {
	var sess token.Session
	err := s.db.WithContext(ctx).First(&sess, "refresh_token = ?", refreshToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, token.ErrNotFound
		}
		return nil, err
	}
	return &sess, nil
}
