package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Jwt      JwtConfig
}
type AppConfig struct {
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JwtConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	return &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "9900"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "postgres"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Jwt: JwtConfig{
			Secret:     getEnv("JWT_SECRET", ""),
			AccessTTL:  parseDuration(getEnv("JWT_ACCESS_TTL", "15m")),
			RefreshTTL: parseDuration(getEnv("JWT_REFRESH_TTL", "24h")),
		},
	}
}

func (c *Config) Validate() {
	if c.Jwt.Secret == "" {
		log.Fatal("JWT secret is required")
	}

	if c.Database.Password == "" && c.App.Env == "production" {
		log.Fatal("DB_PASSWORD is required in production.")
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func parseDuration(s string) time.Duration {
	if len(s) > 0 && s[len(s)-1] == 'd' {
		s = s[:len(s)-1] + "h"
		if d, err := time.ParseDuration(s); err == nil {
			return d * 24
		}
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("Invalid duration value: %s", s)
	}
	return d
}
