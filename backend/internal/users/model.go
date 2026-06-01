package users

import (
	"gorm.io/gorm"

	"github.com/DavoodHakimi/warehouse-app/internal/company"
)

// defining DB Types

type User struct {
	gorm.Model
	FullName    string `gorm:"not null" json:"full_name"`
	UserName    string `gorm:"not null;uniqueIndex" json:"user_name"`
	UserTypeID  uint   `json:"use_type_id"`
	Password    string `gorm:"not null" json:"-"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Email       string `gorm:"uniqueIndex" json:"email"`
	CompanyID   uint   `json:"company_id"`

	// Relationships
	UserType UserType        `gorm:"foreignKey:UserTypeID" json:"user_type,omitempty"`
	Company  company.Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

type UserType struct {
	gorm.Model
	Name string `gorm:"unique; not null" json:"name"`
}
