package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/DavoodHakimi/warehouse-app/internal/company"
	"github.com/DavoodHakimi/warehouse-app/internal/users"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SignUp(r *SignUpRequest) error {
	newCompany := company.Company{
		Name: r.CompanyName,
	}

	hashedPassword, err := users.HashPassword(r.Password)
	if err != nil {
		return err
	}

	newUser := users.User{
		FullName:    r.FullName,
		UserName:    r.UserName,
		UserTypeID:  1,
		Password:    string(hashedPassword),
		PhoneNumber: r.PhoneNumber,
		Email:       r.Email,
	}

	err = s.repo.createComapny(&newCompany, &newUser)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(r *LogInRequest) (string, error, int) {

	user, err := s.repo.readUser(r)
	if err != nil {
		return "", err, http.StatusNotFound
	}
	if v := users.CheckPassword(r.Password, user.Password); v {
		token, err := generateToken(user)
		if err == nil {
			return token, err, http.StatusOK
		}
	}
	return "", errors.New("User can not be recognized"), http.StatusUnauthorized
}

func generateToken(u *users.User) (string, error) {

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    u.ID,
		"username":   u.UserName,
		"company_id": u.CompanyID,
		"role":       u.UserTypeID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	return tokenString, err
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}
