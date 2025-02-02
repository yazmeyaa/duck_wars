package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/yazmeyaa/duck_wars/internal/config"
	"github.com/yazmeyaa/duck_wars/internal/database"
	"github.com/yazmeyaa/duck_wars/internal/game_data"
	"github.com/yazmeyaa/duck_wars/internal/redis"
	character_repository "github.com/yazmeyaa/duck_wars/internal/repository/character"
	user_repository "github.com/yazmeyaa/duck_wars/internal/repository/user"
	authorization_service "github.com/yazmeyaa/duck_wars/internal/service/authorization"
	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
	"github.com/yazmeyaa/duck_wars/internal/transport/transport_http/http_controllers"
	"github.com/yazmeyaa/duck_wars/internal/transport/transport_http/http_router"
	"github.com/yazmeyaa/duck_wars/pkg/logger"
)

var (
	resourcesFilePath = filepath.Join("resources")
)

func main() {
	logger := logger.New()
	cfg := config.New()
	gd := game_data.NewGameData(logger)
	if err := gd.ReadFolderDefaultFiles(resourcesFilePath); err != nil {
		logger.Error("Error while reading game data files.", "error", err.Error())
		os.Exit(1)
		return
	}
	rdb, err := redis.NewRedisClient(cfg)
	if err != nil {
		panic(err)
	}
	db, err := database.New(cfg)
	if err != nil {
		panic(err)
	}
	_, err = redis.NewRedisClient(cfg)
	if err != nil {
		panic(err)
	}
	characterRepository := character_repository.NewPostgres(db)
	userRepository := user_repository.New(db, rdb, gd)
	userService := user_service.New(logger, userRepository, characterRepository)
	characterService := character_service.New(characterRepository, logger)
	characterController := http_controllers.NewCharacterController(userService, characterService)
	authorizationService := authorization_service.New(cfg)
	authController := http_controllers.NewAuthenticationController(authorizationService, userService, logger)
	router := http_router.New(http_router.RouterConstructor{
		CharacterController:      characterController,
		AuthenticationController: authController,
		Logger:                   logger,
		AuthorizationService:     authorizationService,
	})

	if err := http.ListenAndServe(":8000", router); err != nil {
		panic(err)
	}

}
