package sqlconnect

import (
	"github.com/Milua25/go-job-application-tracker/internal/user"
	"gorm.io/gorm"
)

type PostgresStore struct {
	DB   *gorm.DB
	User user.Repository
}

func NewPostgresStore(db *gorm.DB) *PostgresStore {
	return &PostgresStore{
		User: &UserStore{db: db},
		DB:   db,
	}
}
