package db

import (
	"log"
	"os"
	"time"

	"example.com/m/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDB(cfg config.Config) *gorm.DB {
	databaseURL := cfg.DATABASE_URL

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	if databaseURL != "" {
		log.Println("connecting to render database")

		database, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			log.Fatal("failed to connect render database: ", err)
		}

		return database
	}

	log.Println("connecting to local database")

	dsn := "host=localhost user=postgres password=postgres dbname=gorm_demo port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("failed to connect local database: ", err)
	}

	return database
}
