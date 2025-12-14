package database

import (
    "gorm.io/gorm"
    "github.com/glebarez/sqlite"
)

func Connect(path string) (*gorm.DB, error) {
    return gorm.Open(sqlite.Open(path), &gorm.Config{})
}
