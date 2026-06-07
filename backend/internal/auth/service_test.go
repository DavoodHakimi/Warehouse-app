package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_SignUp(t *testing.T) {
	db := newTestDB(t)
	svc := NewService(NewRepository(db))

	req := &SignUpRequest{
		CompanyName:          "Test Company",
		FullName:             "John Doe",
		UserName:             "johndoe",
		Password:             "password123",
		PasswordConfirmation: "password123",
		PhoneNumber:          "09123456789",
		Email:                "john@example.com",
	}

	err := svc.SignUp(req)
	require.NoError(t, err)
}

func TestService_Login_success(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")

	db := newTestDB(t)
	svc := NewService(NewRepository(db))

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

	token, err, status := svc.Login(&LogInRequest{
		UserName: "johndoe",
		Password: "password123",
	})

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.NotEmpty(t, token)
}

func TestService_Login_userNotFound(t *testing.T) {
	db := newTestDB(t)
	svc := NewService(NewRepository(db))

	_, err, status := svc.Login(&LogInRequest{
		UserName: "nobody",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, status)
}

func TestService_Login_wrongPassword(t *testing.T) {
	db := newTestDB(t)
	svc := NewService(NewRepository(db))

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

	_, err, status := svc.Login(&LogInRequest{
		UserName: "johndoe",
		Password: "wrongpassword",
	})

	assert.Error(t, err)
	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestValidateToken_valid(t *testing.T) {
	t.Setenv("JWT_SECRET", "validate-token-secret")

	db := newTestDB(t)
	svc := NewService(NewRepository(db))

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

	token, err, status := svc.Login(&LogInRequest{
		UserName: "johndoe",
		Password: "password123",
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	claims, err := ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, "johndoe", claims["username"])
}

func TestValidateToken_invalid(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	_, err := ValidateToken("not.a.valid.token")
	assert.Error(t, err)
}
