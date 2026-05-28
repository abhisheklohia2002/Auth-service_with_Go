package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	PrivateKeyPath string
	PublicKeyPath  string

	JWTIssuer            string
	JWTKEYID             string
	REFRESH_TOKEN_SECRET string

	AccessTokenExpiryMins   int
	RefreshTokenExpiryHours int

	Host      string
	DBPort    string
	User      string
	Password  string
	Database  string
	DBSSLMode string
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return intValue
}

func LoadDotenv() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables/defaults")
	}

	return Config{
		Port: getEnv("PORT", "5500"),

		PrivateKeyPath: getEnv("PRIVATE_KEY_PATH", "./openSSL/access_private.pem"),
		PublicKeyPath:  getEnv("PUBLIC_KEY_PATH", "./openSSL/public.pem"),

		JWTIssuer:               getEnv("JWT_ISSUER", "http://localhost:5500"),
		JWTKEYID:                getEnv("JWTKEYID", "****"),
		REFRESH_TOKEN_SECRET:    getEnv("REFRESH_TOKEN_SECRET", "******"),
		AccessTokenExpiryMins:   getEnvAsInt("ACCESS_TOKEN_EXPIRY_MINUTES", 60),
		RefreshTokenExpiryHours: getEnvAsInt("REFRESH_TOKEN_EXPIRY_HOURS", 8760),

		Host:      getEnv("DB_HOST", "*****"),
		DBPort:    getEnv("DB_PORT", "*****"),
		User:      getEnv("DB_USER", "*****"),
		Password:  getEnv("DB_PASSWORD", "*****"),
		Database:  getEnv("DB_NAME", "******"),
		DBSSLMode: getEnv("DB_SSLMODE", "disable"),
	}
}
