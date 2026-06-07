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

	c := &company.Company{Name: "Acme Corp"}
	err := repo.createComapny(c)

	require.NoError(t, err)
	assert.NotZero(t, c.ID)

	var stored company.Company
	require.NoError(t, db.First(&stored, c.ID).Error)
	assert.Equal(t, "Acme Corp", stored.Name)
}

func TestRepository_createUser(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)

	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, repo.createComapny(c))

	hashed, err := users.HashPassword("password123")
	require.NoError(t, err)

	u := &users.User{
		FullName:    "Jane Doe",
		UserName:    "janedoe",
		UserTypeID:  1,
		Password:    hashed,
		PhoneNumber: "09123456789",
		Email:       "jane@example.com",
		CompanyID:   c.ID,
	}
	err = repo.createUser(u)

	require.NoError(t, err)
	assert.NotZero(t, u.ID)

	var stored users.User
	require.NoError(t, db.First(&stored, u.ID).Error)
	assert.Equal(t, "janedoe", stored.UserName)
}

func TestRepository_readUser_found(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)

	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, repo.createComapny(c))

	hashed, err := users.HashPassword("password123")
	require.NoError(t, err)

	u := &users.User{
		FullName:    "Jane Doe",
		UserName:    "janedoe",
		UserTypeID:  1,
		Password:    hashed,
		PhoneNumber: "09123456789",
		Email:       "jane@example.com",
		CompanyID:   c.ID,
	}
	require.NoError(t, repo.createUser(u))

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
