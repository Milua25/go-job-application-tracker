package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email           string    `gorm:"uniqueIndex;not null;size:254"`
	PasswordHash    string    `gorm:"not null"`
	FirstName       string    `gorm:"not null;size:50"`
	LastName        string    `gorm:"not null;size:50"`
	Timezone        string    `gorm:"not null;default:UTC"`
	IsActive        bool      `gorm:"not null;default:true"`
	EmailVerifiedAt *time.Time
	IsAdmin         bool `gorm:"not null;default:false"`
	CreatedAt       time.Time
	LastLoginAt     *time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	// Sessions        []token.Session `gorm:"foreignKey:UserID"`
}
