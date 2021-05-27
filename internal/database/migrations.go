package database

import (
	"github.com/dionaditya/go-production-ready-api/internal/comment"
	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(&comment.Comment{})
}
