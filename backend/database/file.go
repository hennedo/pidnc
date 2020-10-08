package database

import (
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Type string `json:"type"`
	Name string `json:"name"`
	Size int64 `json:"size"`
}
