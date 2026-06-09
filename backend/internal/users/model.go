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
	UserTypeID  uint   `json:"user_type_id"`
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
	Name        string `gorm:"unique; not null" json:"name"`
	PersianName string `gorm:"unique" json:"persian_name"`
	Description string `json:"description,omitempty"`
}

type Permission struct {
	gorm.Model
	UserTypeID       uint `json:"user_type_id"`
	PermissionTypeID uint `json:"permission_type_id"`

	// Relationships
	UserType       UserType       `gorm:"foreignKey:UserTypeID" json:"user_type,omitempty"`
	PermissionType PermissionType `gorm:"foreignKey:PermissionTypeID" json:"permission_type,omitempty"`
}

type PermissionType struct {
	gorm.Model
	Name string `gorm:"not null;unique" json:"name"`
}
