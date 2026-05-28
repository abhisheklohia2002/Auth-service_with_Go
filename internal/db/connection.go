package db

import (
	"fmt"
	"log"
	"time"

	"example.com/m/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	env := config.LoadDotenv()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		env.HOST,
		env.USER,
		env.PASSWORD,
		env.DATABASE,
		env.DBPORT,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("failed to get sql db:", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("failed to ping database:", err)
	}

	fmt.Println("Database connected successfully")

	return database
}
