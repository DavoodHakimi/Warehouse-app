package company

import (
	"gorm.io/gorm"
)

// defining DB Types
type Company struct {
	gorm.Model
	Name string `gorm:"unique; not null" json:"name"`
}
