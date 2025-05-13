package tests

import (
	"log"

	"gorm.io/gorm"
)

func CleanupDatabase(db *gorm.DB) {
	if err := db.Exec("DELETE FROM products").Error; err != nil {
		log.Fatalf("Failed to clean up products table: %v", err)
	}

	if err := db.Exec("DELETE FROM users").Error; err != nil {
		log.Fatalf("Failed to clean up users table: %v", err)
	}
}
