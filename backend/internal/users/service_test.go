package users

import (
	"strconv"
	"testing"

	"github.com/DavoodHakimi/warehouse-app/internal/audit"
	"github.com/DavoodHakimi/warehouse-app/internal/company"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupService(t *testing.T) (*Service, *gorm.DB) {
	t.Helper()
	db := newTestDB(t)
	return NewService(NewRepository(db)), db
}

func TestHashPassword_and_CheckPassword(t *testing.T) {
	hash, err := HashPassword("password123")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.True(t, CheckPassword("password123", hash))
	assert.False(t, CheckPassword("wrongpassword", hash))
}

func TestService_AllUsers(t *testing.T) {
	svc, db := setupService(t)
	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, db.Create(c).Error)
	seedUser(t, db, c.ID, "userone", "one@example.com")

	result, err := svc.AllUsers(int(c.ID))
	require.NoError(t, err)
	assert.NotEmpty(t, result.Users)
	assert.Equal(t, "userone", result.Users[0].UserName)
}

func TestService_ReadUser_found(t *testing.T) {
	svc, db := setupService(t)
	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, db.Create(c).Error)
	user := seedUser(t, db, c.ID, "janedoe", "jane@example.com")

	info, err := svc.ReadUser(strconv.Itoa(int(user.ID)))
	require.NoError(t, err)
	assert.Equal(t, "janedoe", info.UserName)
	assert.Equal(t, "jane@example.com", info.Email)
}

func TestService_ReadUser_notFound(t *testing.T) {
	svc, _ := setupService(t)

	info, err := svc.ReadUser("999")
	assert.NoError(t, err)
	assert.Empty(t, info.UserName)
}

func TestService_CreateUser(t *testing.T) {
	svc, db := setupService(t)
	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, db.Create(c).Error)

	err := svc.CreateUser(&CreateUserRequest{
		FullName:             "New User",
		UserName:             "newuser",
		UserTypeID:           1,
		Password:             "password123",
		PasswordConfirmation: "password123",
		PhoneNumber:          "09111111111",
		Email:                "new@example.com",
	}, int(c.ID))
	require.NoError(t, err)

	var count int64
	db.Model(&User{}).Where("user_name = ?", "newuser").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestService_UpdateUser(t *testing.T) {
	svc, db := setupService(t)
	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, db.Create(c).Error)
	user := seedUser(t, db, c.ID, "janedoe", "jane@example.com")

	err := svc.UpdateUser(&UpdateUserRequest{
		ID:          int(user.ID),
		FullName:    "Jane Updated",
		UserName:    "janedoe",
		UserTypeID:  1,
		PhoneNumber: "09123456789",
		Email:       "updated@example.com",
	}, int(user.ID))
	require.NoError(t, err)

	var stored User
	require.NoError(t, db.First(&stored, user.ID).Error)
	assert.Equal(t, "Jane Updated", stored.FullName)
	assert.Equal(t, "updated@example.com", stored.Email)

	var logCount int64
	db.Model(&audit.Log{}).Count(&logCount)
	assert.Greater(t, logCount, int64(0))
}

func TestService_UpdateUser_noChanges(t *testing.T) {
	svc, db := setupService(t)
	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, db.Create(c).Error)
	user := seedUser(t, db, c.ID, "janedoe", "jane@example.com")

	err := svc.UpdateUser(&UpdateUserRequest{
		ID:          int(user.ID),
		FullName:    "Jane Doe",
		UserName:    "janedoe",
		UserTypeID:  1,
		PhoneNumber: "09123456789",
		Email:       "jane@example.com",
	}, int(user.ID))
	assert.Error(t, err)
	assert.Equal(t, "no changes detected", err.Error())
}

func TestService_DeleteUser(t *testing.T) {
	svc, db := setupService(t)
	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, db.Create(c).Error)
	user := seedUser(t, db, c.ID, "janedoe", "jane@example.com")

	err := svc.DeleteUser(int(user.ID))
	require.NoError(t, err)

	var count int64
	db.Model(&User{}).Where("id = ?", user.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestService_modifiedFields(t *testing.T) {
	svc, db := setupService(t)
	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, db.Create(c).Error)
	user := seedUser(t, db, c.ID, "janedoe", "jane@example.com")

	changes := svc.modifiedFields(&UpdateUserRequest{
		ID:          int(user.ID),
		FullName:    "Jane Updated",
		UserName:    "janedoe",
		UserTypeID:  2,
		PhoneNumber: "09123456789",
		Email:       "jane@example.com",
	})

	assert.Contains(t, changes, "FullName")
	assert.Contains(t, changes, "UserTypeID")
	assert.NotContains(t, changes, "Email")
}

