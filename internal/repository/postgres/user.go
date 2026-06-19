package postgres

// Implements the user.Repository interface

import (
	"database/sql"
)

type UserStore struct {
	db *sql.DB // Your actual DB client lives here
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

// GetByID implements user.Repository
// func (s *UserStore) GetByID(ctx context.Context, id string) (*user.User, error) {
// 	// 1. Run your SQL query here...
// 	// 2. Scan data into local variables or DB row structs
// 	// 3. Return a pure *user.User entity
// 	return &user.User{ID: id, Name: "John Doe"}, nil
// }
