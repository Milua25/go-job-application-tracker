package sqlconnect

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Milua25/go-job-application-tracker/internal/user"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func (s *UserStore) GetByID(ctx context.Context, id string) (*user.User, error) {
	slog.Debug("fetching user by id", "user_id", id)
	var u user.User

	err := s.db.WithContext(ctx).First(&u, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Debug("user not found", "user_id", id)
			return nil, user.ErrNotFound
		}
		slog.Error("failed to fetch user by id", "user_id", id, "error", err)
		return nil, err
	}

	slog.Debug("user fetched successfully", "user_id", id, "email", u.Email)
	return &u, nil
}

func (s *UserStore) GetAll(ctx context.Context) ([]*user.User, error) {
	slog.Debug("fetching all users")
	var users []*user.User
	err := s.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		slog.Error("failed to fetch all users", "error", err)
		return nil, err
	}
	slog.Debug("all users fetched successfully", "count", len(users))
	return users, nil
}

func (s *UserStore) Create(ctx context.Context, u *user.User) error {
	slog.Debug("creating user", "email", u.Email, "user_id", u.ID.String())
	return withTx(s.db, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(u).Error; err != nil {
			slog.Error("failed to create user", "email", u.Email, "error", err)
			return err
		}
		slog.Debug("user created successfully", "email", u.Email, "user_id", u.ID.String())
		return nil
	}, &sql.TxOptions{})
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	slog.Debug("fetching user by email", "email", email)
	var u user.User
	err := s.db.WithContext(ctx).First(&u, "email = ?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Debug("user not found by email", "email", email)
			return nil, user.ErrNotFound
		}
		slog.Error("failed to fetch user by email", "email", email, "error", err)
		return nil, err
	}

	slog.Debug("user fetched successfully by email", "email", email, "user_id", u.ID.String())
	return &u, nil
}

func (s *UserStore) Update(ctx context.Context, u *user.User) error {
	slog.Debug("updating user", "user_id", u.ID.String(), "email", u.Email)
	return withTx(s.db, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Save(u).Error; err != nil {
			slog.Error("failed to update user", "user_id", u.ID.String(), "error", err)
			return err
		}
		slog.Debug("user updated successfully", "user_id", u.ID.String())
		return nil
	}, &sql.TxOptions{})
}

func (s *UserStore) Delete(ctx context.Context, id string) error {
	slog.Debug("deleting user", "user_id", id)
	return withTx(s.db, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Delete(&user.User{}, "id = ?", id).Error; err != nil {
			slog.Error("failed to delete user", "user_id", id, "error", err)
			return err
		}
		slog.Debug("user deleted successfully", "user_id", id)
		return nil
	}, &sql.TxOptions{})
}
