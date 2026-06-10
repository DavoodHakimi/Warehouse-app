package main

import (
	"context"
	"log"

	"github.com/DavoodHakimi/warehouse-app/internal/database"
	"github.com/DavoodHakimi/warehouse-app/internal/database/seed"
	"github.com/DavoodHakimi/warehouse-app/internal/router"
)

func main() {

	db := database.GetDB()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Connected to database successfully!")

	seeder := seed.NewSeeder(db)
	if err := seeder.Run(context.Background()); err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	r := router.Setup(db)
	r.Run(":8080")
}
