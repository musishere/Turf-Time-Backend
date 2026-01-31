package models

import "github.com/google/uuid"

type Sports struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	IconUrl    string    `gorm:"column:icon_url;type:varchar(255);not null" json:"iconUrl"`
	MinPlayers int       `gorm:"column:min_players" json:"minPlayers"`
	MaxPlayers int       `gorm:"column:max_players" json:"maxPlayers"`
}
