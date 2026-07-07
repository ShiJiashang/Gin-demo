package config

import (
	"os"
	"time"
)

type Config struct {
	Port      string
	DBPath    string
	JWTSecret string
	JWTExpire time.Duration
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "gin-demo.db"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret"
	}

	return Config{
		Port:      port,
		DBPath:    dbPath,
		JWTSecret: jwtSecret,
		JWTExpire: 2 * time.Hour,
	}
}
