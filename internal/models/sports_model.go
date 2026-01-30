package models

import "github.com/google/uuid"

type Sports struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	IconUrl    string    `gorm:"type:varchar(255);not null" json:"iconUrl"`
	MinPlayers int       `json:"minPlayers"`
	MaxPlayers int       `json:"MaxPlayers"`
}
