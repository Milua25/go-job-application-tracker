package sqlconnect

import (
	"context"

	"github.com/Milua25/go-job-application-tracker/internal/user"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func (s *UserStore) GetByID(ctx context.Context, id string) (*user.User, error) {
	return nil, nil
}

func (s *UserStore) GetAll(ctx context.Context) ([]*user.User, error) {
	return nil, nil
}

func (s *UserStore) Create(ctx context.Context, u *user.User) error {
	return nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	// 	  if err != nil {
	//     return nil, err
	// }
	// if lastLogin.Valid {
	//     user.LastLogin = &lastLogin.Time
	// }
	return nil, nil
}
