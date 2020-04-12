package server

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct{}

func (server *Server) Run() {
	router := echo.New()

	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	if err := router.Start("0.0.0.0:5000"); err != nil {
		log.Fatal(err)
	}
}
