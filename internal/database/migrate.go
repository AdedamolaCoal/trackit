package database

import (
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, models ...interface{}) {
	log.Println("Migrating models...")

	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration completed successfully")
}
