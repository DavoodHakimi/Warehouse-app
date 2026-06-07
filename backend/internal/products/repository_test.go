package products

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_FindByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	c := seedCompany(t, db, "Acme Corp")
	product := seedProduct(t, db, c.ID, "Widget", "PRD-001", 9.99)

	found, err := repo.FindByID("PRD-001")
	require.NoError(t, err)
	assert.Equal(t, product.Name, found.Name)
	assert.Equal(t, "PRD-001", found.ProductNumber)
}

func TestRepository_FindByID_notFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)

	_, err := repo.FindByID("PRD-MISSING")
	assert.Error(t, err)
}

func TestRepository_Create(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	c := seedCompany(t, db, "Acme Corp")

	product := &Product{
		Name:          "New Widget",
		ProductNumber: "PRD-NEW",
		CompanyID:     c.ID,
		DefaultPrice:  12.50,
	}

	err := repo.Create(product)
	require.NoError(t, err)
	assert.NotZero(t, product.ID)
}

func TestRepository_Update(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	c := seedCompany(t, db, "Acme Corp")
	product := seedProduct(t, db, c.ID, "Widget", "PRD-001", 9.99)

	product.Name = "Updated Widget"
	product.DefaultPrice = 19.99
	err := repo.Update(product)
	require.NoError(t, err)

	var stored Product
	require.NoError(t, db.First(&stored, product.ID).Error)
	assert.Equal(t, "Updated Widget", stored.Name)
	assert.Equal(t, 19.99, stored.DefaultPrice)
}

func TestRepository_Delete(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	c := seedCompany(t, db, "Acme Corp")
	product := seedProduct(t, db, c.ID, "Widget", "PRD-001", 9.99)

	err := repo.Delete(product)
	require.NoError(t, err)

	var count int64
	db.Model(&Product{}).Where("id = ?", product.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestRepository_ReadCompanyProducts(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	c := seedCompany(t, db, "Acme Corp")
	other := seedCompany(t, db, "Other Corp")
	seedProduct(t, db, c.ID, "Widget A", "PRD-A", 5.00)
	seedProduct(t, db, c.ID, "Widget B", "PRD-B", 7.00)
	seedProduct(t, db, other.ID, "Other Widget", "PRD-C", 3.00)

	products, err := repo.ReadCompanyProducts(int(c.ID))
	require.NoError(t, err)
	assert.Len(t, products, 2)
}
