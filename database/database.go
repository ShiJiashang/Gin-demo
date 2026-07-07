package database

import (
	"gin_gorm_demo/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(dbPath string) error {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		return err
	}

	DB = db
	return nil
}
