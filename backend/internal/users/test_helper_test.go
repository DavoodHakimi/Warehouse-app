package users

import (
	"testing"

	"github.com/DavoodHakimi/warehouse-app/internal/audit"
	"github.com/DavoodHakimi/warehouse-app/internal/company"
	"github.com/DavoodHakimi/warehouse-app/internal/testutil"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	return testutil.NewTestDB(t, &company.Company{}, &User{}, &audit.Log{})
}
