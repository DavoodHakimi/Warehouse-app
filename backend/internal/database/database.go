package database

import (
	"fmt"
	"os"
	"sync"

	"gorm.io/gorm"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"

	"github.com/DavoodHakimi/warehouse-app/internal/company"
	"github.com/DavoodHakimi/warehouse-app/internal/orders"
	"github.com/DavoodHakimi/warehouse-app/internal/partners"
	"github.com/DavoodHakimi/warehouse-app/internal/products"
	"github.com/DavoodHakimi/warehouse-app/internal/users"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		initDB()
	})
	return db
}

func initDB() {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbHost := os.Getenv("DBHOST")
	dbPort := os.Getenv("DBPORT")
	dbName := os.Getenv("DBNAME")

	dbConn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, dbUser, dbPass, dbName, dbPort)

	connection, err := gorm.Open(postgres.Open(dbConn), &gorm.Config{})

	if err != nil {
		panic("Database connection failed: " + err.Error())
	}
	db = connection
	fmt.Println("Sucessfully connected to DB!")
}

func RunMigrations() error {
	dbObj := GetDB()

	err := dbObj.AutoMigrate(
		&company.Company{},
		&users.User{},
		&users.UserType{},
		&users.Permission{},
		&users.PermissionType{},
		&partners.BusinessPartner{},
		&partners.BusinessPartnerType{},
		&products.Product{},
		&products.Stock{},
		&orders.Currency{},
		&orders.Order{},
		&orders.OrderItem{})
	if err != nil {
		panic("DB migration failed: " + err.Error())
	}

	fmt.Println("Database migration completed successfully!")
	return nil
}
