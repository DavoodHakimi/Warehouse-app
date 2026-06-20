package partners

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ReadCompanyPartners(companyID int) ([]BusinessPartner, error) {

	var partners []BusinessPartner
	res := r.db.Where("company_id = ?", companyID).Find(&partners)
	if res.Error != nil {

		return nil, res.Error
	}
	return partners, nil
}

func (r *Repository) FindByID(id int, companyID int) (*BusinessPartner, error) {
	var partner BusinessPartner
	res := r.db.Joins("BusinessPartnerType").Where("company_id = ?", companyID).First(&partner, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &partner, nil
}

func (r *Repository) Create(partner *BusinessPartner) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(partner).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *Repository) Update(partner *BusinessPartner, companyID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&BusinessPartner{}).Where("id = ? AND company_id = ?", partner.ID, companyID).Updates(map[string]interface{}{
			"name":                     partner.Name,
			"email":                    partner.Email,
			"phone_number":             partner.PhoneNumber,
			"business_partner_type_id":  partner.BusinessPartnerTypeID,
			"contact_name":             partner.ContactName,
			"contact_phone_number":     partner.ContactPhoneNumber,
		}).Error
	})
}

func (r *Repository) Delete(partner *BusinessPartner, companyID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("company_id = ?", companyID).Delete(partner).Error; err != nil {
			return err
		}

		return nil
	})
}
