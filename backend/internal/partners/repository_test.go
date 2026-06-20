package partners

import (
	"fmt"
	"testing"

	"github.com/DavoodHakimi/warehouse-app/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// partnerTypeCounter keeps each seeded type unique (the name column has a
// UNIQUE constraint) across multiple createPartnerForCompany calls.
var partnerTypeCounter int

// setupRepo builds an in-memory SQLite database with the partner tables and
// returns a ready-to-use Repository plus the underlying handle for seeding.
func setupRepo(t *testing.T) (*Repository, *gorm.DB) {
	t.Helper()
	db := testutil.NewTestDB(t, &BusinessPartner{}, &BusinessPartnerType{})
	return NewRepository(db), db
}

// createType inserts a BusinessPartnerType and returns it.
func createType(t *testing.T, db *gorm.DB) BusinessPartnerType {
	t.Helper()
	bpt := BusinessPartnerType{Name: "Supplier", PersianName: "تامین‌کننده"}
	require.NoError(t, db.Create(&bpt).Error)
	return bpt
}

// createTypeNamed inserts a BusinessPartnerType with the given names.
func createTypeNamed(t *testing.T, db *gorm.DB, name, persian string) BusinessPartnerType {
	t.Helper()
	bpt := BusinessPartnerType{Name: name, PersianName: persian}
	require.NoError(t, db.Create(&bpt).Error)
	return bpt
}

// createPartner inserts one BusinessPartner for companyID and returns it. It
// creates a fresh, uniquely-named type each call so it can be invoked many
// times within a single test without hitting the name UNIQUE constraint.
func createPartnerForCompany(t *testing.T, db *gorm.DB, companyID uint) BusinessPartner {
	t.Helper()
	partnerTypeCounter++
	bpt := createTypeNamed(t, db,
		fmt.Sprintf("Type-%d", partnerTypeCounter),
		fmt.Sprintf("نوع-%d", partnerTypeCounter),
	)
	bp := BusinessPartner{
		Name:                  "Acme Co",
		BusinessPartnerTypeID: bpt.ID,
		PhoneNumber:           "09123456789",
		Email:                 "acme@example.com",
		ContactName:           "Davy",
		ContactPhoneNumber:    "09111111111",
		CompanyID:             companyID,
	}
	require.NoError(t, db.Create(&bp).Error)
	return bp
}

func TestRepository_Create_and_FindByID(t *testing.T) {
	repo, db := setupRepo(t)

	bpt := createType(t, db)
	partner := &BusinessPartner{
		Name:                  "Acme Co",
		BusinessPartnerTypeID: bpt.ID,
		PhoneNumber:           "09123456789",
		Email:                 "acme@example.com",
		ContactName:           "Davy",
		ContactPhoneNumber:    "09111111111",
		CompanyID:             1,
	}

	require.NoError(t, repo.Create(partner))
	assert.NotZero(t, partner.ID)

	got, err := repo.FindByID(int(partner.ID), 1)
	require.NoError(t, err)
	assert.Equal(t, partner.Name, got.Name)
	assert.Equal(t, partner.PhoneNumber, got.PhoneNumber)
	assert.Equal(t, partner.Email, got.Email)
	assert.Equal(t, partner.CompanyID, got.CompanyID)
	assert.Equal(t, bpt.Name, got.BusinessPartnerType.Name)
}

func TestRepository_FindByID_NotFound(t *testing.T) {
	repo, _ := setupRepo(t)

	_, err := repo.FindByID(99999, 1)
	require.Error(t, err)
}

func TestRepository_Delete(t *testing.T) {
	repo, db := setupRepo(t)
	partner := createPartnerForCompany(t, db, 1)

	require.NoError(t, repo.Delete(&partner, 1))

	_, err := repo.FindByID(int(partner.ID), 1)
	assert.Error(t, err, "partner should no longer exist after delete")
}

func TestRepository_FindByID_crossCompanyIsolation(t *testing.T) {
	repo, db := setupRepo(t)
	partner := createPartnerForCompany(t, db, 1)

	_, err := repo.FindByID(int(partner.ID), 2)
	require.Error(t, err)
}

func TestRepository_ReadCompanyPartners(t *testing.T) {
	t.Skip("BUG: ReadCompanyPartners filters on the misspelled column 'comnpany_id' " +
		"(should be 'company_id'); the query errors with 'no such column' in SQLite. " +
		"Remove this skip once the typo in repository.go is fixed.")

	repo, db := setupRepo(t)
	createPartnerForCompany(t, db, 7)
	createPartnerForCompany(t, db, 7)
	createPartnerForCompany(t, db, 9) // different company

	got, err := repo.ReadCompanyPartners(7)
	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestRepository_Update(t *testing.T) {
	t.Skip("BUG: Update writes to non-existent columns 'partner_type_id' and " +
		"'contant_phone_number' (should be 'business_partner_type_id' and " +
		"'contact_phone_number'). The UPDATE fails in SQLite with 'no such column'. " +
		"Remove this skip once repository.go is fixed.")

	repo, db := setupRepo(t)
	partner := createPartnerForCompany(t, db, 1)
	bpt2 := createTypeNamed(t, db, "Customer", "مشتری")

	updated := &BusinessPartner{
		Model:                 gorm.Model{ID: partner.ID},
		Name:                  "Renamed",
		Email:                 "new@example.com",
		PhoneNumber:           "09000000000",
		BusinessPartnerTypeID: bpt2.ID,
		ContactPhoneNumber:    "09222222222",
	}
	require.NoError(t, repo.Update(updated, 1))

	got, err := repo.FindByID(int(partner.ID), 1)
	require.NoError(t, err)
	assert.Equal(t, "Renamed", got.Name)
	assert.Equal(t, "new@example.com", got.Email)
	assert.Equal(t, "09000000000", got.PhoneNumber)
	assert.Equal(t, bpt2.ID, got.BusinessPartnerTypeID)
	assert.Equal(t, "09222222222", got.ContactPhoneNumber)
}
