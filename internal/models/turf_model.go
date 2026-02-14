package models

import (
	"time"

	"github.com/google/uuid"
)

type Turf struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	StartTime  int       `gorm:"type:int;not null" json:"startTime"`
	EndTime    int       `gorm:"type:int;not null" json:"endTime"`
	Status     string    `gorm:"type:varchar(50);default:'active'" json:"status"`
	NoOfFields int       `gorm:"type:int;not null" json:"noOfFields"`
	Address    string    `gorm:"type:varchar(255);not null;default:''" json:"address"`
	TurfImages []string  `gorm:"column:turf_images;type:jsonb;serializer:json;not null" json:"turfImages"`
	Longitude  float64   `gorm:"type:double precision;not null;default:0" json:"longitude"`
	Latitude   float64   `gorm:"type:double precision;not null;default:0" json:"latitude"`

	// Relationship: Turf belongs to a User (Owner/Admin)
	OwnerID uuid.UUID `gorm:"type:uuid;not null" json:"ownerId"`
	Owner   User      `gorm:"foreignKey:OwnerID;references:ID" json:"owner,omitempty"`
	// Standard timestamps
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
