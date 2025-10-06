package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"todo-backend/models"
)

func New(path string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("file:%s?_foreign_keys=on", path)
	gormDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := gormDB.AutoMigrate(&models.Todo{}); err != nil {
		return nil, err
	}
	return gormDB, nil
}
