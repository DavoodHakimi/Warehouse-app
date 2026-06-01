package main

import (
	"log"

	"github.com/DavoodHakimi/warehouse-app/internal/database"
	"github.com/DavoodHakimi/warehouse-app/internal/router"
)

func main() {

	// dbConn := database.GetDB()

	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Connected to database successfully!")

	r := router.Setup()
	r.Run(":8080")
}
