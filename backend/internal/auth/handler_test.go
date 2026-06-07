package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHandler_SignUpHandler_success(t *testing.T) {
	db := newTestDB(t)
	handler := NewHandler(NewService(NewRepository(db)))

	body, err := json.Marshal(SignUpRequest{
		CompanyName:          "Test Company",
		FullName:             "John Doe",
		UserName:             "johndoe",
		Password:             "password123",
		PasswordConfirmation: "password123",
		PhoneNumber:          "09123456789",
		Email:                "john@example.com",
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.SignUpHandler(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "User created successfully", resp["message"])
}

func TestHandler_LogInHandler_success(t *testing.T) {
	t.Setenv("JWT_SECRET", "handler-test-secret")

	db := newTestDB(t)
	svc := NewService(NewRepository(db))
	handler := NewHandler(svc)

	signUpReq := &SignUpRequest{
		CompanyName:          "Test Company",
		FullName:             "John Doe",
		UserName:             "johndoe",
		Password:             "password123",
		PasswordConfirmation: "password123",
		PhoneNumber:          "09123456789",
		Email:                "john@example.com",
	}
	require.NoError(t, svc.SignUp(signUpReq))

	body, err := json.Marshal(LogInRequest{
		UserName: "johndoe",
		Password: "password123",
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.LogInHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp["token"])
}

func TestHandler_LogInHandler_userNotFound(t *testing.T) {
	db := newTestDB(t)
	handler := NewHandler(NewService(NewRepository(db)))

	body, err := json.Marshal(LogInRequest{
		UserName: "nobody",
		Password: "password123",
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.LogInHandler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandler_LogInHandler_wrongPassword(t *testing.T) {
	db := newTestDB(t)
	svc := NewService(NewRepository(db))
	handler := NewHandler(svc)

	signUpReq := &SignUpRequest{
		CompanyName:          "Test Company",
		FullName:             "John Doe",
		UserName:             "johndoe",
		Password:             "password123",
		PasswordConfirmation: "password123",
		PhoneNumber:          "09123456789",
		Email:                "john@example.com",
	}
	require.NoError(t, svc.SignUp(signUpReq))

	body, err := json.Marshal(LogInRequest{
		UserName: "johndoe",
		Password: "wrongpassword",
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.LogInHandler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMeHandler_success(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", 1)
	c.Set("company_id", 2)
	c.Set("username", "johndoe")

	MeHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp meResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 1, resp.UserID)
	assert.Equal(t, 2, resp.CompanyID)
	assert.Equal(t, "johndoe", resp.UserName)
}

func TestMeHandler_missingContextKeys(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	MeHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
