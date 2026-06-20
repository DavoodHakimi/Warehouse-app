package products

import (
	"testing"

	"github.com/DavoodHakimi/warehouse-app/internal/audit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupService(t *testing.T) (*Service, *gorm.DB) {
	t.Helper()
	db := newTestDB(t)
	return NewService(NewRepository(db)), db
}

func TestService_AllProducts(t *testing.T) {
	svc, db := setupService(t)
	c := seedCompany(t, db, "Acme Corp")
	seedProduct(t, db, c.ID, "Widget A", "PRD-A", 5.00)

	result, err := svc.AllProducts(int(c.ID))
	require.NoError(t, err)
	assert.Len(t, result.Products, 1)
	assert.Equal(t, "Widget A", result.Products[0].Name)
}

func TestService_ReadProduct_found(t *testing.T) {
	svc, db := setupService(t)
	c := seedCompany(t, db, "Acme Corp")
	seedProduct(t, db, c.ID, "Widget", "PRD-001", 9.99)

	info, err := svc.ReadProduct("PRD-001", int(c.ID))
	require.NoError(t, err)
	assert.Equal(t, "Widget", info.Name)
	assert.Equal(t, "PRD-001", info.ProductNumber)
}

func TestService_ReadProduct_notFound(t *testing.T) {
	svc, db := setupService(t)
	c := seedCompany(t, db, "Acme Corp")

	_, err := svc.ReadProduct("PRD-MISSING", int(c.ID))
	assert.Error(t, err)
}

func TestService_CreateProduct(t *testing.T) {
	svc, db := setupService(t)
	c := seedCompany(t, db, "Acme Corp")

	err := svc.CreateProduct(&ProductRequest{
		Name:         "New Widget",
		IsFrozen:     false,
		DefaultPrice: 15.00,
	}, int(c.ID))
	require.NoError(t, err)

	var count int64
	db.Model(&Product{}).Where("name = ?", "New Widget").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestService_UpdateProduct(t *testing.T) {
	svc, db := setupService(t)
	c := seedCompany(t, db, "Acme Corp")
	product := seedProduct(t, db, c.ID, "Widget", "PRD-001", 9.99)

	err := svc.UpdateProduct(&UpdateProductRequest{
		ID:            int(product.ID),
		Name:          "Updated Widget",
		ProductNumber: "PRD-001",
		IsFrozen:      true,
		DefaultPrice:  19.99,
	}, 1, int(c.ID))
	require.NoError(t, err)

	var stored Product
	require.NoError(t, db.First(&stored, product.ID).Error)
	assert.Equal(t, "Updated Widget", stored.Name)
	assert.True(t, stored.IsFrozen)
	assert.Equal(t, 19.99, stored.DefaultPrice)

	var logCount int64
	db.Model(&audit.Log{}).Count(&logCount)
	assert.Greater(t, logCount, int64(0))
}

func TestService_UpdateProduct_noChanges(t *testing.T) {
	svc, db := setupService(t)
	c := seedCompany(t, db, "Acme Corp")
	product := seedProduct(t, db, c.ID, "Widget", "PRD-001", 9.99)

	err := svc.UpdateProduct(&UpdateProductRequest{
		ID:            int(product.ID),
		Name:          "Widget",
		ProductNumber: "PRD-001",
		IsFrozen:      false,
		DefaultPrice:  9.99,
	}, 1, int(c.ID))
	assert.Error(t, err)
	assert.Equal(t, "no changes detected", err.Error())
}

func TestService_DeleteProduct(t *testing.T) {
	svc, db := setupService(t)
	c := seedCompany(t, db, "Acme Corp")
	seedProduct(t, db, c.ID, "Widget", "PRD-001", 9.99)

	err := svc.DeleteProduct("PRD-001", int(c.ID))
	require.NoError(t, err)

	var count int64
	db.Model(&Product{}).Where("product_number = ?", "PRD-001").Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestService_modifiedFields(t *testing.T) {
	svc, db := setupService(t)
	c := seedCompany(t, db, "Acme Corp")
	product := seedProduct(t, db, c.ID, "Widget", "PRD-001", 9.99)

	changes := svc.modifiedFields(&UpdateProductRequest{
		ID:            int(product.ID),
		Name:          "Updated Widget",
		ProductNumber: "PRD-001",
		IsFrozen:      true,
		DefaultPrice:  9.99,
	}, int(c.ID))

	assert.Contains(t, changes, "Name")
	assert.Contains(t, changes, "IsFrozen")
	assert.NotContains(t, changes, "DefaultPrice")
}
