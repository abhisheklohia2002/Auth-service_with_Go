package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL != "" {
		database, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect render database: ", err)
		}

		return database
	}

	dsn := "host=localhost user=postgres password=postgres dbname=gorm_demo port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect local database: ", err)
	}

	return database
}