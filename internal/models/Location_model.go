package models

type Location struct {
	ID        string  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Latitude  float64 `gorm:"type:decimal(10,8);not null" json:"latitude"`
	Longitude float64 `gorm:"type:decimal(11,8);not null" json:"longitude"`
}
