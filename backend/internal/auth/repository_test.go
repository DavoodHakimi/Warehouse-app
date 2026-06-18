package auth

import (
	"testing"

	"github.com/DavoodHakimi/warehouse-app/internal/company"
	"github.com/DavoodHakimi/warehouse-app/internal/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_createComapny(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)

	hashed, err := users.HashPassword("password123")
	require.NoError(t, err)

	c := &company.Company{Name: "Acme Corp"}
	u := &users.User{
		FullName:    "Jane Doe",
		UserName:    "janedoe",
		UserTypeID:  1,
		Password:    hashed,
		PhoneNumber: "09123456789",
		Email:       "jane@example.com",
	}
	err = repo.createComapny(c, u)

	require.NoError(t, err)
	assert.NotZero(t, c.ID)
	assert.NotZero(t, u.ID)
	assert.Equal(t, c.ID, u.CompanyID, "user must be linked to the created company")

	var storedCompany company.Company
	require.NoError(t, db.First(&storedCompany, c.ID).Error)
	assert.Equal(t, "Acme Corp", storedCompany.Name)

	var storedUser users.User
	require.NoError(t, db.First(&storedUser, u.ID).Error)
	assert.Equal(t, "janedoe", storedUser.UserName)
	assert.Equal(t, c.ID, storedUser.CompanyID)
}

// TestRepository_createComapny_rollbackOnUserError proves the transaction is
// atomic: when user creation fails, the already-created company is rolled back.
func TestRepository_createComapny_rollbackOnUserError(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)

	hashed, err := users.HashPassword("password123")
	require.NoError(t, err)

	// First signup succeeds.
	c1 := &company.Company{Name: "First Corp"}
	u1 := &users.User{
		FullName:   "Jane Doe",
		UserName:   "janedoe",
		UserTypeID: 1,
		Password:   hashed,
		Email:      "jane@example.com",
	}
	require.NoError(t, repo.createComapny(c1, u1))

	// Second signup reuses u1's UserName, so user creation fails on the
	// unique index. The new company must be rolled back.
	c2 := &company.Company{Name: "Second Corp"}
	u2 := &users.User{
		FullName:   "John Doe",
		UserName:   "janedoe", // duplicate -> constraint violation
		UserTypeID: 1,
		Password:   hashed,
		Email:      "john@example.com",
	}
	err = repo.createComapny(c2, u2)
	require.Error(t, err)

	var companyCount int64
	require.NoError(t, db.Model(&company.Company{}).Count(&companyCount).Error)
	assert.Equal(t, int64(1), companyCount, "second company should have been rolled back")

	var userCount int64
	require.NoError(t, db.Model(&users.User{}).Count(&userCount).Error)
	assert.Equal(t, int64(1), userCount, "second user should not have been persisted")
}

func TestRepository_readUser_found(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)

	hashed, err := users.HashPassword("password123")
	require.NoError(t, err)

	c := &company.Company{Name: "Acme Corp"}
	u := &users.User{
		FullName:    "Jane Doe",
		UserName:    "janedoe",
		UserTypeID:  1,
		Password:    hashed,
		PhoneNumber: "09123456789",
		Email:       "jane@example.com",
	}
	require.NoError(t, repo.createComapny(c, u))

	found, err := repo.readUser(&LogInRequest{
		UserName: "janedoe",
		Password: "password123",
	})

	require.NoError(t, err)
	assert.Equal(t, "janedoe", found.UserName)
	assert.Equal(t, hashed, found.Password)
}

func TestRepository_readUser_notFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)

	_, err := repo.readUser(&LogInRequest{UserName: "missing"})

	assert.Error(t, err)
}
