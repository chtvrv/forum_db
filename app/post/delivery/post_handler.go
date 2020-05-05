package delivery

import (
	"net/http"

	"github.com/chtvrv/forum_db/app/middleware"
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/pkg/errors"

	"github.com/chtvrv/forum_db/app/post"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type PostHandler struct {
	Usecase post.Usecase
}

func CreateHandler(router *echo.Echo, usecase post.Usecase) {
	handler := &PostHandler{
		Usecase: usecase,
	}

	router.POST("/api/thread/:slug_or_id/create", handler.Create, middleware.ReadBody, middleware.ReadThreadIdentifier, middleware.Headers)
}

func (postHandler *PostHandler) Create(ctx echo.Context) error {
	var posts models.Posts
	postsBody := ctx.Get("body").([]byte)
	err := posts.UnmarshalJSON(postsBody)
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	slugOrID := ctx.Get("SlugOrID").(string)

	err, msg := postHandler.Usecase.Create(&posts, slugOrID)
	// Успешно создали
	if err == nil {
		response, err := posts.MarshalJSON()
		if err != nil {
			log.Error(err)
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.String(http.StatusCreated, string(response))
	}
	// // Конфликт
	if err == errors.ErrConflict {
		//previousThread, _ := threadHandler.Usecase.GetThreadBySlug(thread.Slug)
		// response, err := previousThread.MarshalJSON()
		// if err != nil {
		// 	log.Error(err)
		// 	return ctx.NoContent(http.StatusInternalServerError)
		// }
		return ctx.JSON(http.StatusConflict, *msg)
	}
	// // Незарегистрированная ошибка
	log.Error(err)
	return ctx.JSON(errors.ResolveErrorToCode(err), *msg)
}
