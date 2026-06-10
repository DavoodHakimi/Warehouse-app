package partners

import (
	"gorm.io/gorm"

	"github.com/DavoodHakimi/warehouse-app/internal/company"
)

// defining DB Types
type BusinessPartner struct {
	gorm.Model
	Name                  string `gorm:"not null" json:"name"`
	BusinessPartnerTypeID uint   `json:"business_partner_type_id"`
	PhoneNumber           string `gorm:"not null" json:"phone_number"`
	Email                 string `json:"email,omitempty"`
	ContactName           string `json:"contact_name,omitempty"`
	ContactPhoneNumber    string `json:"contact_phone_number,omitempty"`
	CompanyID             uint   `json:"company_id"`

	// Relationships
	BusinessPartnerType BusinessPartnerType `gorm:"foreignKey:BusinessPartnerTypeID" json:"business_partner_type,omitempty"`
	Company             company.Company     `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

type BusinessPartnerType struct {
	gorm.Model
	Name        string `gorm:"unique; not null" json:"name"`
	PersianName string `gorm:"unique" json:"persian_name"`
}
