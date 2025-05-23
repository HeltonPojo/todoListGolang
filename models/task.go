package models

import (
	"gorm.io/gorm"
)

type Tasks struct {
	gorm.Model
	Title       string `gorm:"not null" json:"title"`
	Description string `json:"description"`
	Status      string `gorm:"default:'peding'" json:"status"`
	UserID      uint   `json:"user_id"`
}
