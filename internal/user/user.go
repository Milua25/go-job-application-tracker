package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID  `db:"id"`
	Email           string     `db:"email"`
	PasswordHash    string     `db:"password_hash"`
	FirstName       string     `db:"first_name"`
	LastName        string     `db:"last_name"`
	Timezone        string     `db:"timezone"`
	IsActive        bool       `db:"is_active,default:true"`
	EmailVerifiedAt *time.Time `db:"email_verified_at"`
	Token           *string    `db:"token"`
	RefreshToken    *string    `db:"refresh_token"`
	IsAdmin         bool       `db:"is_admin,default:false"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}
