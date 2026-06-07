package products

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupHandler(t *testing.T) (*Handler, *gorm.DB) {
	t.Helper()
	db := newTestDB(t)
	return NewHandler(NewService(NewRepository(db))), db
}

func TestHandler_AllProductsHandler_forbidden(t *testing.T) {
	handler, db := setupHandler(t)
	c := seedCompany(t, db, "Acme Corp")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/products", nil)
	ctx.Set("role", 2)
	ctx.Set("company_id", int(c.ID))

	handler.AllProductsHandler(ctx)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestHandler_AllProductsHandler_success(t *testing.T) {
	handler, db := setupHandler(t)
	c := seedCompany(t, db, "Acme Corp")
	seedProduct(t, db, c.ID, "Widget", "PRD-001", 9.99)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/products", nil)
	ctx.Set("role", 1)
	ctx.Set("company_id", int(c.ID))

	handler.AllProductsHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ProductsInfo
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp.Products, 1)
}

func TestHandler_ProductInfoHandler_success(t *testing.T) {
	handler, db := setupHandler(t)
	c := seedCompany(t, db, "Acme Corp")
	seedProduct(t, db, c.ID, "Widget", "PRD-001", 9.99)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/products/PRD-001", nil)
	ctx.Params = gin.Params{{Key: "productNumber", Value: "PRD-001"}}

	handler.ProductInfoHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ProductInfoResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "Widget", resp.Name)
	assert.Equal(t, "PRD-001", resp.ProductNumber)
}

func TestHandler_ProductCreationHandler_success(t *testing.T) {
	handler, db := setupHandler(t)
	c := seedCompany(t, db, "Acme Corp")

	body, err := json.Marshal(ProductRequest{
		Name:         "New Widget",
		IsFrozen:     false,
		DefaultPrice: 12.50,
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("company_id", int(c.ID))

	handler.ProductCreationHandler(ctx)

	assert.Equal(t, http.StatusCreated, w.Code)

	var count int64
	db.Model(&Product{}).Where("name = ?", "New Widget").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestHandler_ProductUpdateHandler_success(t *testing.T) {
	handler, db := setupHandler(t)
	c := seedCompany(t, db, "Acme Corp")
	product := seedProduct(t, db, c.ID, "Widget", "PRD-001", 9.99)

	body, err := json.Marshal(UpdateProductRequest{
		ID:            int(product.ID),
		Name:          "Updated Widget",
		ProductNumber: "PRD-001",
		IsFrozen:      true,
		DefaultPrice:  19.99,
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPatch, "/products/PRD-001", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("user_id", 1)

	handler.ProductUpdateHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var stored Product
	require.NoError(t, db.First(&stored, product.ID).Error)
	assert.Equal(t, "Updated Widget", stored.Name)
	assert.True(t, stored.IsFrozen)
	assert.Equal(t, 19.99, stored.DefaultPrice)
}

func TestHandler_ProductUpdateHandler_unauthorized(t *testing.T) {
	handler, db := setupHandler(t)
	c := seedCompany(t, db, "Acme Corp")
	product := seedProduct(t, db, c.ID, "Widget", "PRD-001", 9.99)

	body, err := json.Marshal(UpdateProductRequest{
		ID:            int(product.ID),
		Name:          "Updated Widget",
		ProductNumber: "PRD-001",
		IsFrozen:      true,
		DefaultPrice:  19.99,
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPatch, "/products/PRD-001", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.ProductUpdateHandler(ctx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandler_ProductDeleteHandler_success(t *testing.T) {
	handler, db := setupHandler(t)
	c := seedCompany(t, db, "Acme Corp")
	seedProduct(t, db, c.ID, "Widget", "PRD-001", 9.99)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/products/PRD-001", nil)
	ctx.Params = gin.Params{{Key: "productNumber", Value: "PRD-001"}}

	handler.ProductDeleteHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var count int64
	db.Model(&Product{}).Where("product_number = ?", "PRD-001").Count(&count)
	assert.Equal(t, int64(0), count)
}
