package repository

import (
	"gorm.io/gorm"

	"github.com/DavoodHakimi/warehouse-app/internal/audit"
)

type Repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) RecordAudit(log *audit.Log) {
	audit.Record(r.db, log)
}

func (r *Repository[T]) ReadAll(companyID int) ([]T, error) {

	var allObjects []T
	res := r.db.Where("company_id = ?", companyID).Find(&allObjects)
	if res.Error != nil {

		return nil, res.Error
	}
	return allObjects, nil
}

func (r *Repository[T]) FindByID(id int) (*T, error) {
	var entity T
	res := r.db.First(&entity, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &entity, nil
}

func (r *Repository[T]) Create(entity *T) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(entity).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *Repository[T]) Update(id uint, updates map[string]interface{}) error {
	var entity T
	return r.db.Model(&entity).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository[T]) Delete(entity *T) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(entity).Error; err != nil {
			return err
		}

		return nil
	})
}
