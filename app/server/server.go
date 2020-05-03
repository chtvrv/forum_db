package server

import (
	"log"

	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	userHandler "github.com/chtvrv/forum_db/app/user/delivery"
	userRepository "github.com/chtvrv/forum_db/app/user/repository"
	userUsecase "github.com/chtvrv/forum_db/app/user/usecase"

	forumHandler "github.com/chtvrv/forum_db/app/forum/delivery"
	forumRepository "github.com/chtvrv/forum_db/app/forum/repository"
	forumUsecase "github.com/chtvrv/forum_db/app/forum/usecase"

	config "github.com/chtvrv/forum_db/pkg/config"
)

type Server struct {
	configReader *config.ConfigReader
}

func (server *Server) Run() {
	server.configReader = config.CreateConfigReader()
	router := echo.New()

	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	postgresConn, err := pgx.NewConnPool(server.configReader.GetDBConn())
	if err != nil {
		log.Fatal(err)
	}

	// Пользователь
	uRepository := userRepository.CreateRepository(postgresConn)
	uUsecase := userUsecase.CreateUsecase(uRepository)
	userHandler.CreateHandler(router, uUsecase)

	// Форум
	fRepository := forumRepository.CreateRepository(postgresConn)
	fUsecase := forumUsecase.CreateUsecase(fRepository, uRepository)
	forumHandler.CreateHandler(router, fUsecase)

	if err := router.Start(server.configReader.GetServerConn()); err != nil {
		log.Fatal(err)
	}
}
