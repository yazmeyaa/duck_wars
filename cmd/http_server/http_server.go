package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/yazmeyaa/ton_game/internal/config"
	"github.com/yazmeyaa/ton_game/internal/database"
	"github.com/yazmeyaa/ton_game/internal/game_data"
	"github.com/yazmeyaa/ton_game/internal/redis"
	characterRepository "github.com/yazmeyaa/ton_game/internal/repository/character"
	userRepository "github.com/yazmeyaa/ton_game/internal/repository/user"
	authorizationService "github.com/yazmeyaa/ton_game/internal/service/authorization"
	characterService "github.com/yazmeyaa/ton_game/internal/service/character"
	userService "github.com/yazmeyaa/ton_game/internal/service/user"
	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/http_controllers"
	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/http_router"
	"github.com/yazmeyaa/ton_game/pkg/logger"
)

var (
	resourcesFilePath = filepath.Join("resources")
)

func main() {
	appLogger := logger.New("http_server")
	cfg := config.New()
	gd := game_data.NewGameData(appLogger)
	if err := gd.ReadFolderDefaultFiles(resourcesFilePath); err != nil {
		appLogger.Error("Error while reading game data files.", "error", err.Error())
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
	characterRepo := characterRepository.NewPostgres(db)
	userRepo := userRepository.New(db, rdb, gd)
	userSvc := userService.New(appLogger, userRepo, characterRepo)
	characterSvc := characterService.New(characterRepo, appLogger)
	characterController := http_controllers.NewCharacterController(userSvc, characterSvc, appLogger)
	authSvc := authorizationService.New(cfg)
	authController := http_controllers.NewAuthenticationController(authSvc, userSvc, appLogger)
	resourcesController := http_controllers.NewResourcesController(gd, appLogger)
	router := http_router.New(http_router.RouterConstructor{
		CharacterController:      characterController,
		AuthenticationController: authController,
		Logger:                   appLogger,
		AuthorizationService:     authSvc,
		ResourcesController:      resourcesController,
	})

	if err := http.ListenAndServe(":8000", router); err != nil {
		panic(err)
	}

}
