package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/DavoodHakimi/warehouse-app/internal/company"
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

func TestHandler_AllUsersHandler_forbidden(t *testing.T) {
	handler, db := setupHandler(t)
	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, db.Create(c).Error)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/users", nil)
	ctx.Set("role", 2)
	ctx.Set("company_id", int(c.ID))

	handler.AllUsersHandler(ctx)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestHandler_AllUsersHandler_success(t *testing.T) {
	handler, db := setupHandler(t)
	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, db.Create(c).Error)
	seedUser(t, db, c.ID, "userone", "one@example.com")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/users", nil)
	ctx.Set("role", 1)
	ctx.Set("company_id", int(c.ID))

	handler.AllUsersHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp UsersInfo
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp.Users)
}

func TestHandler_UserInfoHandler_success(t *testing.T) {
	handler, db := setupHandler(t)
	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, db.Create(c).Error)
	user := seedUser(t, db, c.ID, "janedoe", "jane@example.com")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/users/"+strconv.Itoa(int(user.ID)), nil)
	ctx.Params = gin.Params{{Key: "userID", Value: strconv.Itoa(int(user.ID))}}

	handler.UserInfoHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp UserInfoResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "janedoe", resp.UserName)
}

func TestHandler_UserCreationHandler_success(t *testing.T) {
	handler, db := setupHandler(t)
	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, db.Create(c).Error)

	body, err := json.Marshal(CreateUserRequest{
		FullName:             "New User",
		UserName:             "newuser",
		UserTypeID:           1,
		Password:             "password123",
		PasswordConfirmation: "password123",
		PhoneNumber:          "09111111111",
		Email:                "new@example.com",
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("company_id", int(c.ID))

	handler.UserCreationHandler(ctx)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHandler_UserUpdateHandler_success(t *testing.T) {
	handler, db := setupHandler(t)
	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, db.Create(c).Error)
	user := seedUser(t, db, c.ID, "janedoe", "jane@example.com")

	body, err := json.Marshal(UpdateUserRequest{
		ID:          int(user.ID),
		FullName:    "Jane Updated",
		UserName:    "janedoe",
		UserTypeID:  1,
		PhoneNumber: "09123456789",
		Email:       "updated@example.com",
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPatch, "/users/"+strconv.Itoa(int(user.ID)), bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Params = gin.Params{{Key: "userID", Value: strconv.Itoa(int(user.ID))}}

	handler.UserUpdateHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var stored User
	require.NoError(t, db.First(&stored, user.ID).Error)
	assert.Equal(t, "Jane Updated", stored.FullName)
	assert.Equal(t, "updated@example.com", stored.Email)
}

func TestHandler_UserDeleteHandler_success(t *testing.T) {
	handler, db := setupHandler(t)
	c := &company.Company{Name: "Acme Corp"}
	require.NoError(t, db.Create(c).Error)
	user := seedUser(t, db, c.ID, "janedoe", "jane@example.com")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/users/"+strconv.Itoa(int(user.ID)), nil)
	ctx.Params = gin.Params{{Key: "userID", Value: strconv.Itoa(int(user.ID))}}

	handler.UserDeleteHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var count int64
	db.Model(&User{}).Where("id = ?", user.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
