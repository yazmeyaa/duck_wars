package main

import (
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"
	"github.com/yazmeyaa/duck_wars/config"
	socket_server "github.com/yazmeyaa/duck_wars/internal/transport/socket"
	"github.com/yazmeyaa/duck_wars/internal/user"
	"github.com/yazmeyaa/duck_wars/pkg/database"
	"github.com/yazmeyaa/duck_wars/pkg/logger"
)

func main() {
	logg := logger.New()
	config, err := config.New(logg)
	if err != nil {
		log.Fatalf("Cannot load config file. Error message: %s", err.Error())
	}

	db, err := database.NewDB(config, logg)
	if err != nil {
		log.Fatalf("Cannot get database instance. Error message: %s", err.Error())
	}

	if err := runMigrations(db); err != nil {
		log.Fatalf("Error while migrate database: %s", err.Error())
	}

	userRepo := user.NewRepository(logg, db)
	usr, _ := userRepo.GetById(1)
	username := "Testing"
	userRepo.Update(usr, &user.UpdateUserDTO{Username: &username})
	usr, err = userRepo.GetById(1)

	if err != nil {
		logger.LogError(logg, "main", "Error while get user", err)
		return
	}

	logg.Info("Created user", "ID", usr.ID, "CreatedAt", usr.CreatedAt, "telegramId", usr.TelegramID)

	socket_server.NewServer()
}

func runMigrations(db *sql.DB) error {
	migrationsDir := "pkg/database/migrations"

	err := goose.Up(db, migrationsDir)
	if err != nil {
		return err
	}

	log.Println("Migrations successfull completed!")
	return nil
}
