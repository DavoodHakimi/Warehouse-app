package orders

import (
	"fmt"

	"github.com/DavoodHakimi/warehouse-app/internal/audit"
	"github.com/DavoodHakimi/warehouse-app/internal/products"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ReadCompanyOrders(companyID int) ([]Order, error) {

	var orders []Order
	res := r.db.Joins("BusinessPartner").Joins("Currency").Where("orders.company_id = ?", companyID).Find(&orders)
	if res.Error != nil {

		return nil, res.Error
	}
	return orders, nil
}

func (r *Repository) FindByID(orderID uint, companyID uint) (*Order, error) {
	var order Order
	res := r.db.Joins("BusinessPartner").Joins("Currency").Preload("OrderItems").Where("orders.company_id = ?", companyID).First(&order, orderID)
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

func (r *Repository) Update(order *Order, companyID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&Order{}).Where("id = ? and orders.company_id = ?", order.ID, companyID).Updates(map[string]interface{}{
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

func (r *Repository) StatusUpdate(orderID uint, newStatus string, expectedStatus ...string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		query := tx.Model(&Order{}).Where("id = ?", orderID)
		if len(expectedStatus) > 0 && expectedStatus[0] != "" {
			query = query.Where("status = ?", expectedStatus[0])
		}
		result := query.Update("status", newStatus)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("status transition conflict: order %d is not in expected state", orderID)
		}
		return nil
	})
}

func (r *Repository) RecordAudit(log *audit.Log) {
	audit.Record(r.db, log)
}

func (r *Repository) ApproveSaleTransaction(orderID uint, items []OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&Order{}).
			Where("id = ? AND status = ?", orderID, "Pending").
			Update("status", "Approved")
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("status transition conflict: order %d is not Pending", orderID)
		}

		for _, item := range items {
			if err := r.Reserve(tx, item.ProductID, item.Quantity); err != nil {
				return err
			}
		}
		return tx.Model(&Order{}).Where("id = ?", orderID).Update("status", "Approved").Error
	})
}

func (r *Repository) ShipSaleTransaction(orderID uint, items []OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := r.Fulfill(tx, item.ProductID, item.Quantity); err != nil {
				return err
			}
		}
		return tx.Model(&Order{}).Where("id = ?", orderID).Update("status", "Shipped").Error
	})
}

func (r *Repository) ReceivePurchaseTransaction(orderID uint, items []OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := r.ReceiveStock(tx, item.ProductID, item.Quantity); err != nil {
				return err
			}
		}
		return tx.Model(&Order{}).Where("id = ?", orderID).Update("status", "Received").Error
	})
}

func (r *Repository) CancelSaleTransaction(orderID uint, items []OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := r.Release(tx, item.ProductID, item.Quantity); err != nil {
				return err
			}
		}
		return tx.Model(&Order{}).Where("id = ?", orderID).Update("status", "Canceled").Error
	})
}

func (r *Repository) Reserve(tx *gorm.DB, productID uint, amount int) error {
	return adjustStock(tx, productID, -amount, amount)
}

func (r *Repository) Release(tx *gorm.DB, productID uint, amount int) error {
	return adjustStock(tx, productID, amount, -amount)
}

func (r *Repository) Fulfill(tx *gorm.DB, productID uint, amount int) error {
	return adjustStock(tx, productID, 0, -amount)
}

func (r *Repository) ReceiveStock(tx *gorm.DB, productID uint, amount int) error {
	return adjustStock(tx, productID, amount, 0)
}

func adjustStock(tx *gorm.DB, productID uint, deltaAvail, deltaRes int) error {
	result := tx.Model(&products.Stock{}).
		Where("product_id = ?", productID).
		Updates(map[string]interface{}{
			"available_stock": gorm.Expr("available_stock + ?", deltaAvail),
			"reserved_stock":  gorm.Expr("reserved_stock + ?", deltaRes),
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return tx.Create(&products.Stock{
			ProductID:      productID,
			AvailableStock: deltaAvail,
			ReservedStock:  deltaRes,
		}).Error
	}

	return nil
}
