package delivery

import (
	"net/http"

	// "log"
	// "net/http"
	// "strconv"
	// "time"
	"github.com/chtvrv/forum_db/app/middleware"
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/app/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type UserHandler struct {
	Usecase user.Usecase
}

func CreateHandler(router *echo.Echo, usecase user.Usecase) {
	handler := &UserHandler{
		Usecase: usecase,
	}

	router.POST("/user/:nickname/create", handler.Create, middleware.ReadUserNickname, middleware.ReadBody)

}

func (userHandler *UserHandler) Create(ctx echo.Context) error {
	var user models.User
	user.Nickname = ctx.Get("nickname").(string)

	userBody := ctx.Get("body").([]byte)
	err := user.UnmarshalJSON(userBody)

	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	err = userHandler.Usecase.Create(&user)
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusOK)
}
