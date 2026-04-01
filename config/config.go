package config

import (
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	ServerPort string
	RedisAddr  string
}

func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "caretop"),
		DBPassword: getEnv("DB_PASSWORD", "caretop123"),
		DBName:     getEnv("DB_NAME", "caretop"),
		JWTSecret:  getEnv("JWT_SECRET", "caretop-secret-key-change-in-production"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
