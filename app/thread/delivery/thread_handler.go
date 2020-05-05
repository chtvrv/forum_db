package delivery

import (
	"net/http"

	"github.com/chtvrv/forum_db/app/middleware"
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/pkg/errors"

	"github.com/chtvrv/forum_db/app/thread"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type ThreadHandler struct {
	Usecase thread.Usecase
}

func CreateHandler(router *echo.Echo, usecase thread.Usecase) {
	handler := &ThreadHandler{
		Usecase: usecase,
	}

	router.POST("/api/forum/:slug/create", handler.Create, middleware.ReadBody, middleware.ReadForumSlug, middleware.Headers)
}

func (threadHandler *ThreadHandler) Create(ctx echo.Context) error {
	var thread models.Thread
	threadBody := ctx.Get("body").([]byte)
	err := thread.UnmarshalJSON(threadBody)
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	thread.Forum = ctx.Get("slug").(string)

	err, msg := threadHandler.Usecase.Create(&thread)
	// Успешно создали
	if err == nil {
		response, err := thread.MarshalJSON()
		if err != nil {
			log.Error(err)
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.String(http.StatusCreated, string(response))
	}
	// Конфликт
	if err == errors.ErrConflict {
		previousThread, _ := threadHandler.Usecase.GetThreadBySlug(thread.Slug)
		response, err := previousThread.MarshalJSON()
		if err != nil {
			log.Error(err)
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.String(http.StatusConflict, string(response))
	}
	// Незарегистрированная ошибка
	log.Error(err)
	return ctx.JSON(errors.ResolveErrorToCode(err), *msg)
}
