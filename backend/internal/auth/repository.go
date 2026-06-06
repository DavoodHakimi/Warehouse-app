package auth

import (
	"github.com/DavoodHakimi/warehouse-app/internal/company"
	"github.com/DavoodHakimi/warehouse-app/internal/users"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) createComapny(c *company.Company) error {
	r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(c).Error; err != nil {
			return err
		}

		return nil
	})
	return nil
}

func (r *Repository) createUser(u *users.User) error {
	r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(u).Error; err != nil {
			return err
		}

		return nil
	})
	return nil
}

func (r *Repository) readUser(d *LogInRequest) (*users.User, error) {
	var u users.User
	result := r.db.Where("user_name = ?", d.UserName).First(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return &u, nil

}
