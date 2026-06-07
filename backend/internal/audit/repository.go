package audit

import "gorm.io/gorm"

func Record(tx *gorm.DB, log *Log) error {
	return tx.Create(&log).Error
}
