package auth

import (
	"errors"
	"log"

	"github.com/DavoodHakimi/warehouse-app/internal/company"
	"github.com/DavoodHakimi/warehouse-app/internal/database"
	"github.com/DavoodHakimi/warehouse-app/internal/users"
)

func createComapny(c *company.Company) error {
	db := database.GetDB()
	res := db.Create(c)
	if res.Error != nil {

		return res.Error
	}
	return nil
}

func createUser(u *users.User) error {
	db := database.GetDB()

	res := db.Create(&u)
	if res.Error != nil {
		log.Fatal("Error creating new User:", res.Error)
		return errors.New("Failed to create user")
	}
	return nil
}

func readUser(d *LogInRequest) (*users.User, error) {
	db := database.GetDB()
	var u users.User
	result := db.Where("user_name = ?", d.UserName).First(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return &u, nil

}
