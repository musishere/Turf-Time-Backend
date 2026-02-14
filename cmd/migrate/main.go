// migrate runs GORM AutoMigrate to sync the database with your Go models.
// Usage: go run ./cmd/migrate (or: make migrate-up)
//
// When adding a new NOT NULL column to an existing table, add a default in the
// gorm tag so existing rows get a value, e.g.:
//
//	NewField string `gorm:"type:varchar(255);not null;default:''"`
package main

import (
	"log"

	"github.com/musishere/sportsApp/config"
	"github.com/musishere/sportsApp/internal/database"
	"github.com/musishere/sportsApp/internal/models"
)

func main() {
	cfg := config.LoadConfig()
	db := database.ConnectDatabase(cfg)

	log.Println("Running AutoMigrate (syncing schema from Go models)...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.Location{},
		&models.Sports{},
		&models.Turf{},
	); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	log.Println("Database schema up to date.")
}
