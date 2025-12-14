package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"go-backend/internal/models"
)

func Connect(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.Device{}); err != nil {
		return nil, err
	}
	return db, nil
}
