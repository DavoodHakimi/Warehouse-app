package products

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ReadCompanyProducts(companyID int) ([]Product, error) {

	var products []Product
	res := r.db.Where("company_id = ?", companyID).Find(&products)
	if res.Error != nil {

		return nil, res.Error
	}
	return products, nil
}

func (r *Repository) FindByID(pNum string) (*Product, error) {
	var product Product
	res := r.db.Where("product_number = ?", pNum).First(&product)
	if res.Error != nil {
		return nil, res.Error
	}
	return &product, nil
}

func (r *Repository) Create(product *Product) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *Repository) Update(prod *Product) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&Product{}).Where("id = ?", prod.ID).Updates(map[string]interface{}{
			"name":          prod.Name,
			"is_frozen":     prod.IsFrozen,
			"default_price": prod.DefaultPrice,
		}).Error
	})
}

func (r *Repository) Delete(product *Product) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(product).Error; err != nil {
			return err
		}

		return nil
	})
}
