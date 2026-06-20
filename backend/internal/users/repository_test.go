package users

import (
	"testing"

	"github.com/DavoodHakimi/warehouse-app/internal/company"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedCompany(t *testing.T, db *gorm.DB) *company.Company {
	t.Helper()

	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, db.Create(c).Error)
	return c
}

func seedUser(t *testing.T, db *gorm.DB, companyID uint, userName, email string) *User {
	t.Helper()

	hashed, err := HashPassword("password123")
	require.NoError(t, err)

	user := &User{
		FullName:    "Jane Doe",
		UserName:    userName,
		UserTypeID:  1,
		Password:    hashed,
		PhoneNumber: "09123456789",
		Email:       email,
		CompanyID:   companyID,
	}
	require.NoError(t, db.Create(user).Error)
	return user
}

func TestRepository_FindByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	c := seedCompany(t, db)
	user := seedUser(t, db, c.ID, "janedoe", "jane@example.com")

	found, err := repo.FindByID(int(user.ID), int(c.ID))
	require.NoError(t, err)
	assert.Equal(t, "janedoe", found.UserName)
}

func TestRepository_FindByID_notFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	c := seedCompany(t, db)

	_, err := repo.FindByID(999, int(c.ID))
	assert.Error(t, err)
}

func TestRepository_FindByID_crossCompanyIsolation(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	c := seedCompany(t, db)
	user := seedUser(t, db, c.ID, "crossuser", "cross@example.com")

	_, err := repo.FindByID(int(user.ID), 9999)
	assert.Error(t, err)
}

func TestRepository_Create(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	c := seedCompany(t, db)

	hashed, err := HashPassword("password123")
	require.NoError(t, err)

	user := &User{
		FullName:    "New User",
		UserName:    "newuser",
		UserTypeID:  1,
		Password:    hashed,
		PhoneNumber: "09111111111",
		Email:       "new@example.com",
		CompanyID:   c.ID,
	}

	err = repo.Create(user)
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestRepository_Update(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	c := seedCompany(t, db)
	user := seedUser(t, db, c.ID, "janedoe", "jane@example.com")

	user.FullName = "Jane Updated"
	err := repo.Update(user, int(c.ID))
	require.NoError(t, err)

	var stored User
	require.NoError(t, db.First(&stored, user.ID).Error)
	assert.Equal(t, "Jane Updated", stored.FullName)
}

func TestRepository_Delete(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	c := seedCompany(t, db)
	user := seedUser(t, db, c.ID, "janedoe", "jane@example.com")

	err := repo.Delete(user, int(c.ID))
	require.NoError(t, err)

	var count int64
	db.Model(&User{}).Where("id = ?", user.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestRepository_ReadCompanyUsers(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	c := seedCompany(t, db)
	seedUser(t, db, c.ID, "userone", "one@example.com")
	seedUser(t, db, c.ID, "usertwo", "two@example.com")

	users, err := repo.ReadCompanyUsers(int(c.ID))
	require.NoError(t, err)
	assert.NotEmpty(t, users)
}
