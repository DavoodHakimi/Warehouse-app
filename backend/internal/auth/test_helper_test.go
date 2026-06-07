package auth

import (
	"testing"

	"github.com/DavoodHakimi/warehouse-app/internal/company"
	"github.com/DavoodHakimi/warehouse-app/internal/testutil"
	"github.com/DavoodHakimi/warehouse-app/internal/users"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	return testutil.NewTestDB(t, &company.Company{}, &users.User{})
}
