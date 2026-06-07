package audit

import (
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	EntityType string `gorm:"not null"` // "order", "product", "stock"
	EntityID   uint   `gorm:"not null"` // which record changed
	Event      string `gorm:"not null"` // "created", "updated", "deleted"
	Field      string
	OldValue   string
	NewValue   string
	ByUserID   uint `gorm:"not null"`
}
