package audit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestRecord(t *testing.T) {
	db := newTestDB(t)

	log := &Log{
		EntityType: "user",
		EntityID:   1,
		Event:      "updated",
		Field:      "Email",
		OldValue:   "old@example.com",
		NewValue:   "new@example.com",
		ByUserID:   1,
	}

	err := Record(db, log)
	require.NoError(t, err)
	assert.NotZero(t, log.ID)

	var stored Log
	require.NoError(t, db.First(&stored, log.ID).Error)
	assert.Equal(t, "user", stored.EntityType)
	assert.Equal(t, uint(1), stored.EntityID)
	assert.Equal(t, "updated", stored.Event)
	assert.Equal(t, "Email", stored.Field)
	assert.Equal(t, "old@example.com", stored.OldValue)
	assert.Equal(t, "new@example.com", stored.NewValue)
	assert.Equal(t, uint(1), stored.ByUserID)
}

func TestRecord_withinTransaction(t *testing.T) {
	db := newTestDB(t)

	err := db.Transaction(func(tx *gorm.DB) error {
		log := &Log{
			EntityType: "order",
			EntityID:   42,
			Event:      "created",
			ByUserID:   2,
		}
		return Record(tx, log)
	})
	require.NoError(t, err)

	var count int64
	require.NoError(t, db.Model(&Log{}).Count(&count).Error)
	assert.Equal(t, int64(1), count)
}
