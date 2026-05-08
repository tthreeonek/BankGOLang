package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
	JWTSecret  string
	AESKey     string // 32 байта для AES-256
	HMACSecret string
	SMTPHost   string
	SMTPPort   string
	SMTPUser   string
	SMTPPass   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "bankdb"),
		SSLMode:    getEnv("SSL_MODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", "supersecretkey"),
		AESKey:     getEnv("AES_KEY", "12345678901234567890123456789012"), // 32 символа
		HMACSecret: getEnv("HMAC_SECRET", "hmacsecretkey"),
		SMTPHost:   getEnv("SMTP_HOST", "smtp.example.com"),
		SMTPPort:   getEnv("SMTP_PORT", "587"),
		SMTPUser:   getEnv("SMTP_USER", "noreply@example.com"),
		SMTPPass:   getEnv("SMTP_PASS", "password"),
	}, nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
