package orders

import (
	"github.com/DavoodHakimi/warehouse-app/internal/company"
	"github.com/DavoodHakimi/warehouse-app/internal/partners"
	"github.com/DavoodHakimi/warehouse-app/internal/products"
	"gorm.io/gorm"
)

// defining DB Types
type Currency struct {
	gorm.Model
	Name        string `gorm:"unique;not null" json:"name"`
	PersianName string `gorm:"unique;not null" json:"persian_name"`
}

type Order struct {
	gorm.Model
	OrderType         string  `gorm:"not null;check:order_type IN ('sale','purchase')" json:"order_type"`
	OrderNumber       string  `gorm:"uniqueIndex" json:"order_number"`
	CompanyID         uint    `gorm:"not null" json:"company_id"`
	Status            string  `gorm:"not null;check:status IN ('Pending','Approved','Shipped','Packed','Waiting','Received','Canceled')" json:"status"`
	BusinessPartnerID uint    `gorm:"not null" json:"business_partner_id"`
	CurrencyID        uint    `gorm:"not null" json:"currency_id"`
	ExchangeRate      float64 `gorm:"type:decimal(10,4);default:1.0" json:"exchange_rate"`

	// Relationships
	Company         company.Company          `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Currency        Currency                 `gorm:"foreignKey:CurrencyID" json:"currency,omitempty"`
	BusinessPartner partners.BusinessPartner `gorm:"foreignKey:BusinessPartnerID" json:"business_partner,omitempty"`
	OrderItems      []OrderItem              `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
}

type OrderItem struct {
	gorm.Model
	ProductID    uint    `json:"product_id"`
	OrderID      uint    `json:"order_id"`
	Quantity     int     `gorm:"not null" binding:"gte=1" json:"quantity"`
	PerItemPrice float64 `gorm:"type:decimal(25,2)" json:"per_item_price"`

	// Relationships
	Product products.Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Order   Order            `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}
