package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name     string    `gorm:"type:varchar(255);not null" json:"name"`
	Email    string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password string    `gorm:"type:varchar(255);not null" json:"-"`
	Role     string    `gorm:"type:varchar(50);not null;default:player" json:"role"`
	IsActive bool      `gorm:"type:boolean;default:true" json:"is_active"`

	// Relations
	Location Location `gorm:"constraint:OnDelete:CASCADE;" json:"location"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
