package user

// Core Domain Entity
import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID  `db:"id"`
	Email           string     `db:"email"`
	PasswordHash    string     `db:"password_hash"`
	FullName        string     `db:"full_name"`
	Timezone        string     `db:"timezone"`
	IsActive        bool       `db:"is_active"`
	EmailVerifiedAt *time.Time `db:"email_verified_at"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}
