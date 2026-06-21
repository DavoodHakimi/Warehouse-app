package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/DavoodHakimi/warehouse-app/internal/database"
	"github.com/DavoodHakimi/warehouse-app/internal/database/seed"
	"github.com/DavoodHakimi/warehouse-app/internal/router"
)

func main() {
	file, _ := os.OpenFile("../../app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()

	logger := slog.New(slog.NewMultiHandler(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}),
	))

	slog.SetDefault(logger)

	db := database.GetDB()

	if err := database.RunMigrations(db); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database successfully")

	seeder := seed.NewSeeder(db)
	if err := seeder.Run(context.Background()); err != nil {
		slog.Error("seeding operation failed", "error", err)
		os.Exit(1)
	}

	r := router.Setup(db)
	slog.Info("server starting", "port", ":8080")
	r.Run(":8080")
}
