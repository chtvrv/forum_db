package delivery

import (
	"net/http"

	"github.com/chtvrv/forum_db/app/middleware"
	"github.com/chtvrv/forum_db/app/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type ServiceHandler struct {
	Usecase service.Usecase
}

func CreateHandler(router *echo.Echo, usecase service.Usecase) {
	handler := &ServiceHandler{
		Usecase: usecase,
	}

	router.POST("/api/service/clear", handler.ClearDB)
	router.GET("/api/service/status", handler.GetStatus, middleware.Headers)
}

func (serviceHandler *ServiceHandler) GetStatus(ctx echo.Context) error {
	status, err := serviceHandler.Usecase.GetStatus()
	if err != nil {
		log.Error(err)
		ctx.NoContent(http.StatusInternalServerError)
	}
	response, err := status.MarshalJSON()
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(response))

}

func (serviceHandler *ServiceHandler) ClearDB(ctx echo.Context) error {
	err := serviceHandler.Usecase.ClearDB()
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.NoContent(http.StatusOK)
}
