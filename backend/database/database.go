package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabase(path string) *gorm.DB {
	var err error
	db, err = gorm.Open(sqlite.Open(path + "/sqlite.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&File{})
	return db
}

func GetDatabase() *gorm.DB {
	return db
}

