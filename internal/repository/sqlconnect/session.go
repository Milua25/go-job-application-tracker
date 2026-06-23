package sqlconnect

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Milua25/go-job-application-tracker/internal/token"
	"gorm.io/gorm"
)

type SessionStore struct {
	db *gorm.DB
}

func (s *SessionStore) CreateRefreshToken(ctx context.Context, newUserSession *token.Session) error {
	// Implement the logic to create a session in the database
	err := withTx(s.db, func(tx *gorm.DB) error {
		if err := tx.Create(newUserSession).Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{})

	if err != nil {
		return err
	}
	return nil
}

func (s *SessionStore) GetRefreshToken(ctx context.Context, refreshToken string) (*token.Session, error) {
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

func (s *SessionStore) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	slog.Info("deleting session with refresh token", "refresh_token", refreshToken)
	err := withTx(s.db, func(tx *gorm.DB) error {
		if err := tx.Delete(&token.Session{}, "refresh_token = ?", refreshToken).Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{})

	if err != nil {
		return err
	}
	return nil
}

func (s *SessionStore) DeleteSessionsByEmail(ctx context.Context, email string) error {
	slog.Info("deleting sessions for user", "email", email)
	return withTx(s.db, func(tx *gorm.DB) error {
		return tx.Delete(&token.Session{}, "user_email = ?", email).Error
	}, &sql.TxOptions{})
}
