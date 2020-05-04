package delivery

import (
	"net/http"

	"github.com/chtvrv/forum_db/app/middleware"
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/app/user"
	"github.com/chtvrv/forum_db/pkg/errors"
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

	router.POST("/api/user/:nickname/create", handler.Create, middleware.ReadUserNickname, middleware.ReadBody, middleware.Headers)
	router.GET("/api/user/:nickname/profile", handler.GetUser, middleware.ReadUserNickname, middleware.Headers)
	router.POST("/api/user/:nickname/profile", handler.Update, middleware.ReadUserNickname, middleware.ReadBody, middleware.Headers)
}

func (userHandler *UserHandler) Create(ctx echo.Context) error {
	user := models.User{Nickname: ctx.Get("nickname").(string)}
	userBody := ctx.Get("body").([]byte)
	err := user.UnmarshalJSON(userBody)
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	err = userHandler.Usecase.Create(&user)
	// Успешно создали
	if err == nil {
		response, err := user.MarshalJSON()
		if err != nil {
			log.Error(err)
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.String(http.StatusCreated, string(response))
	}
	// Возник конфликт
	if err == errors.ErrConflict {
		var previousUsers models.Users

		previousByNickname, err := userHandler.Usecase.GetUserByNickname(user.Nickname)
		if err != nil && err != errors.ErrNoRows {
			log.Error(err)
			return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		}
		if previousByNickname != nil {
			previousUsers = append(previousUsers, *previousByNickname)
		}

		previousByEmail, err := userHandler.Usecase.GetUserByEmail(user.Email)
		if err != nil && err != errors.ErrNoRows {
			log.Error(err)
			return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		}
		if previousByEmail != nil && (previousByNickname == nil || previousByEmail.Nickname != previousByNickname.Nickname) {
			previousUsers = append(previousUsers, *previousByEmail)
		}

		response, err := previousUsers.MarshalJSON()
		if err != nil {
			log.Error(err)
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.String(http.StatusConflict, string(response))
	}
	// Незарегистрированная ошибка
	log.Error(err)
	return ctx.JSON(errors.ResolveErrorToCode(err), err.Error())
}

func (userHandler *UserHandler) GetUser(ctx echo.Context) error {
	nickname := ctx.Get("nickname").(string)

	user, err := userHandler.Usecase.GetUserByNickname(nickname)
	if err != nil {
		log.Error(err)
		//return ctx.String(errors.ResolveErrorToCode(err), err.Error())
		return ctx.JSON(errors.ResolveErrorToCode(err), errors.CreateMessageNotFoundUser(nickname))
	}

	response, err := user.MarshalJSON()
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(response))
}

func (userHandler *UserHandler) Update(ctx echo.Context) error {
	updatedUser := models.User{Nickname: ctx.Get("nickname").(string)}
	userBody := ctx.Get("body").([]byte)
	err := updatedUser.UnmarshalJSON(userBody)
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	err, msg := userHandler.Usecase.Update(&updatedUser)
	if err != nil {
		log.Error(err)
		//return ctx.JSON(errors.ResolveErrorToCode(err), err.Error())
		return ctx.JSON(errors.ResolveErrorToCode(err), *msg)
	}

	response, err := updatedUser.MarshalJSON()
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(response))
}
