package main

import (
	"net/http"

	"github.com/yazmeyaa/duck_wars/internal/config"
	"github.com/yazmeyaa/duck_wars/internal/database"
	character_repository "github.com/yazmeyaa/duck_wars/internal/repository/character"
	user_repository "github.com/yazmeyaa/duck_wars/internal/repository/user"
	authorization_service "github.com/yazmeyaa/duck_wars/internal/service/authorization"
	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
	"github.com/yazmeyaa/duck_wars/internal/transport/http/http_controllers"
	"github.com/yazmeyaa/duck_wars/internal/transport/http/http_router"
	"github.com/yazmeyaa/duck_wars/pkg/logger"
)

func main() {
	cfg := config.New()
	db, err := database.New(cfg)
	logger := logger.New()
	if err != nil {
		panic(err)
	}
	characterRepository := character_repository.NewPostgres(db)
	userRepository := user_repository.New(db)
	userService := user_service.New(logger, userRepository, characterRepository)
	characterService := character_service.New(characterRepository, logger)
	characterController := http_controllers.NewCharacterController(userService, *characterService)
	authorizationService := authorization_service.New(cfg)
	router := http_router.New(http_router.RouterConstructor{
		CharacterController:  characterController,
		Logger:               logger,
		AuthorizationService: authorizationService,
	})

	if err := http.ListenAndServe(":8000", router); err != nil {
		panic(err)
	}

}
