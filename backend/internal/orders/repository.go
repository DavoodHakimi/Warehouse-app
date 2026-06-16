package orders

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ReadCompanyOrders(companyID int) ([]Order, error) {

	var orders []Order
	res := r.db.Joins("BusinessPartner").Joins("Currency").Where("company_id = ?", companyID).Find(&orders)
	if res.Error != nil {

		return nil, res.Error
	}
	return orders, nil
}

func (r *Repository) FindByID(orderID uint) (*Order, error) {
	var order Order
	res := r.db.Joins("BusinessPartner").Joins("Currency").First(&order, orderID)
	if res.Error != nil {
		return nil, res.Error
	}
	return &order, nil
}

func (r *Repository) Create(order *Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *Repository) Update(order *Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&Order{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
			"order_type":          order.OrderType,
			"status":              order.Status,
			"business_partner_id": order.BusinessPartnerID,
			"currency_id":         order.CurrencyID,
			"exchange_rate":       order.ExchangeRate,
		}).Error
	})
}

func (r *Repository) Delete(order *Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(order).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *Repository) StatusUpdate(orderID uint, newStatus string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&Order{}).Where("id = ?", orderID).Updates(map[string]interface{}{
			"status": newStatus,
		}).Error
	})
}
