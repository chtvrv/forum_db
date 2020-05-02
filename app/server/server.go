package server

import (
	"log"

	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	userHandler "github.com/chtvrv/forum_db/app/user/delivery"
	userRepository "github.com/chtvrv/forum_db/app/user/repository"
	userUsecase "github.com/chtvrv/forum_db/app/user/usecase"
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

	// user
	uRepository := userRepository.CreateRepository(postgresConn)
	uUsecase := userUsecase.CreateUsecase(uRepository)
	userHandler.CreateHandler(router, uUsecase)

	if err := router.Start(server.configReader.GetServerConn()); err != nil {
		log.Fatal(err)
	}
}
