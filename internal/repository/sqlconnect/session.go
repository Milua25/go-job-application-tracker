package sqlconnect

import (
	"context"
	"database/sql"

	"github.com/Milua25/go-job-application-tracker/internal/session"
	"gorm.io/gorm"
)

type SessionStore struct {
	db *gorm.DB
}

func (s *SessionStore) CreateRefreshToken(ctx context.Context, newUserSession *session.Session) error {
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

func (s *SessionStore) GetRefreshToken(ctx context.Context, refreshToken string) (*session.Session, error) {
	var sess session.Session
	err := s.db.WithContext(ctx).First(&sess, "refresh_token = ?", refreshToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, session.ErrNotFound
		}
		return nil, err
	}
	return &sess, nil
}

func (s *SessionStore) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	err := withTx(s.db, func(tx *gorm.DB) error {
		if err := tx.Delete(&session.Session{}, "refresh_token = ?", refreshToken).Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{})

	if err != nil {
		return err
	}
	return nil
}
