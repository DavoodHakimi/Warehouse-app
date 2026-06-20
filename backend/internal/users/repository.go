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
	res := r.db.Where("company_id = ?", companyID).Find(&users)
	if res.Error != nil {

		return nil, res.Error
	}
	return users, nil
}

func (r *Repository) FindByID(id int, companyID int) (*User, error) {
	var user User
	res := r.db.Where("company_id = ?", companyID).First(&user, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (r *Repository) Create(user *User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *Repository) Update(user *User, companyID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&User{}).Where("id = ?", user.ID).Where("company_id = ?", companyID).Updates(map[string]interface{}{
			"full_name":    user.FullName,
			"user_name":    user.UserName,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"user_type_id": user.UserTypeID,
		}).Error
	})
}

func (r *Repository) Delete(user *User, companyID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("company_id = ?", companyID).Delete(user).Error; err != nil {
			return err
		}

		return nil
	})
}
