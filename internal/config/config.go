package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
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

	HOST     string
	USER     string
	PASSWORD string
	DATABASE string
	PORT     string
	DBPORT   string

	DATABASE_URL    string
	JWT_PRIVATE_KEY string
	JWT_PUBLIC_KEY  string

	CLOUDINARY_CLOUD_NAME string
	CLOUDINARY_API_KEY    string
	CLOUDINARY_API_SECRET string

	APP_ENV string
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func normalizePEM(value string) string {
	return strings.ReplaceAll(value, `\n`, "\n")
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

func NewCloudinary() *cloudinary.Cloudinary {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Fatal("failed to initialize cloudinary: ", err)
	}

	return cld
}
func LoadDotenv() Config {
	wd, _ := os.Getwd()
	log.Println("working dir:", wd)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("failed to load .env:", err)
	}
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables/defaults")
	}

	return Config{
		Port: getEnv("PORT", "5500"),

		PrivateKeyPath: getEnv("PRIVATE_KEY_PATH", "./openSSL/access_private.pem"),
		PublicKeyPath:  getEnv("PUBLIC_KEY_PATH", "./openSSL/public.pem"),

		JWTIssuer:               getEnv("JWT_ISSUER", "http://localhost:5500"),
		JWTKEYID:                getEnv("JWTKEYID", "****"),
		REFRESH_TOKEN_SECRET:    getEnv("REFRESH_TOKEN_SECRET", ""),
		AccessTokenExpiryMins:   getEnvAsInt("ACCESS_TOKEN_EXPIRY_MINUTES", 60),
		RefreshTokenExpiryHours: getEnvAsInt("REFRESH_TOKEN_EXPIRY_HOURS", 8760),

		HOST:     getEnv("DB_HOST", ""),
		USER:     getEnv("DB_USER", ""),
		PASSWORD: getEnv("DB_PASSWORD", ""),
		DATABASE: getEnv("DB_NAME", ""),
		DBPORT:   getEnv("DB_PORT", ""),

		DATABASE_URL:    getEnv("DATABASE_URL", ""),
		JWT_PRIVATE_KEY: normalizePEM(getEnv("JWT_PRIVATE_KEY", "")),
		JWT_PUBLIC_KEY:  normalizePEM(getEnv("JWT_PUBLIC_KEY", "")),

		CLOUDINARY_CLOUD_NAME: getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CLOUDINARY_API_KEY:    getEnv("CLOUDINARY_API_KEY", ""),
		CLOUDINARY_API_SECRET: getEnv("CLOUDINARY_API_SECRET", ""),

		APP_ENV: getEnv("APP_ENV", ""),
	}
}
