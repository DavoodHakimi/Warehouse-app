package users

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ReadCompanyUsers(companyID int) ([]User, error) {

	var users []User
	res := r.db.Find(&users).Where("comnpany_id = ?", companyID)
	if res.Error != nil {

		return nil, res.Error
	}
	return users, nil
}

func (r *Repository) FindByID(id int) (*User, error) {
	var user User
	res := r.db.First(&user, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (r *Repository) Create(user *User) error {
	r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		return nil
	})

	return nil
}

func (r *Repository) Update(user *User) error {
	r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(user).Error; err != nil {
			return err
		}

		return nil
	})
	return nil
}

func (r *Repository) Delete(user *User) error {
	r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(user).Error; err != nil {
			return err
		}

		return nil
	})
	return nil
}
