package repositories

import "gorm.io/gorm"

type TurfRepostitory struct {
	db *gorm.DB
}

func NewTurfRepository(db *gorm.DB) *TurfRepostitory {
	return &TurfRepostitory{
		db: db,
	}
}
