package products

import (
	"gorm.io/gorm"

	"github.com/DavoodHakimi/warehouse-app/internal/company"
)

// defining DB Types

type Product struct {
	gorm.Model
	Name          string  `gorm:"not null" json:"name"`
	ProductNumber string  `gorm:"unique; not null" json:"product_number"`
	CompanyID     uint    `json:"company_id"`
	IsFrozen      bool    `gorm:"default:false" json:"is_frozen"`
	DefaultPrice  float64 `gorm:"type:decimal(25,2)" json:"default_price"`

	// Relationships
	Company company.Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

type Stock struct {
	gorm.Model
	ProductID      uint `json:"product_id"`
	AvailableStock int  `gorm:"not null;default:0;check:available_stock >= 0" binding:"gte=0" json:"available_stock"`
	ReservedStock  int  `gorm:"not null;default:0;check:reserved_stock >= 0" binding:"gte=0" json:"reserved_stock"`

	// Relationships
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
