package products

import (
	"testing"

	"github.com/DavoodHakimi/warehouse-app/internal/audit"
	"github.com/DavoodHakimi/warehouse-app/internal/company"
	"github.com/DavoodHakimi/warehouse-app/internal/testutil"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	return testutil.NewTestDB(t, &company.Company{}, &Product{}, &Stock{}, &audit.Log{})
}

func seedCompany(t *testing.T, db *gorm.DB, name string) *company.Company {
	t.Helper()

	c := &company.Company{Name: name}
	require.NoError(t, db.Create(c).Error)
	return c
}

func seedProduct(t *testing.T, db *gorm.DB, companyID uint, name, productNumber string, price float64) *Product {
	t.Helper()

	p := &Product{
		Name:          name,
		ProductNumber: productNumber,
		CompanyID:     companyID,
		DefaultPrice:  price,
	}
	require.NoError(t, db.Create(p).Error)
	return p
}
