package database

import (
	"fmt"
	"go-task/config"
	"go-task/models"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	p := config.GetEnv("DB_PORT")

	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable timezone=Asia/Jakarta",
		config.GetEnv("DB_HOST"),
		config.GetEnv("DB_USER"),
		config.GetEnv("DB_PASSWORD"),
		config.GetEnv("DB_NAME"),
		port,
	)

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Failed to connect the database: %v", err))
	}

	fmt.Println("Database connection open.")
	DB.AutoMigrate(&models.Users{})
	DB.AutoMigrate(&models.Products{})
	fmt.Println("Database migrated success.")
}
