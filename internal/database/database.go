package database

import (
	"expense-tracker/internal/config"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func Connect(cfg *config.Config) *Database {
	dsn := buildDsn(cfg.Database)

	gormConfig := &gorm.Config{
		Logger: buildLogger(cfg.App.Env),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	sqlDb.SetMaxOpenConns(25) // in prod this will be 100
	sqlDb.SetMaxIdleConns(20)
	sqlDb.SetConnMaxLifetime(5 * time.Minute)

	log.Printf("Connected to database")
	return &Database{db}
}

func buildDsn(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)
}

func buildLogger(env string) logger.Interface {
	if env == "development" {
		return logger.Default.LogMode(logger.Info)
	}
	return logger.Default.LogMode(logger.Error)
}
