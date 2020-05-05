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
	router.POST("api/thread/:slug_or_id/vote", handler.VoteForThread, middleware.ReadBody, middleware.ReadThreadIdentifier, middleware.Headers)
	router.GET("api/thread/:slug_or_id/details", handler.GetThread, middleware.ReadThreadIdentifier, middleware.Headers)
	router.POST("api/thread/:slug_or_id/details", handler.UpdateThread, middleware.ReadBody, middleware.ReadThreadIdentifier, middleware.Headers)
	router.GET("api/thread/:slug_or_id/posts", handler.GetPosts, middleware.ReadThreadIdentifier, middleware.ReadGetPostsQuery, middleware.Headers)
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

func (threadHandler *ThreadHandler) VoteForThread(ctx echo.Context) error {
	var vote models.Vote
	voteBody := ctx.Get("body").([]byte)
	err := vote.UnmarshalJSON(voteBody)
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	slugOrID := ctx.Get("SlugOrID").(string)

	thread, err, msg := threadHandler.Usecase.VoteForThread(&vote, slugOrID)
	if err != nil {
		log.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), *msg)
	}

	response, err := thread.MarshalJSON()
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(response))
}

func (threadHandler *ThreadHandler) GetThread(ctx echo.Context) error {
	slugOrID := ctx.Get("SlugOrID").(string)

	thread, err, msg := threadHandler.Usecase.GetThreadByIdentifier(slugOrID)
	if err != nil {
		log.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), *msg)
	}
	response, err := thread.MarshalJSON()
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(response))
}

func (threadHandler *ThreadHandler) GetPosts(ctx echo.Context) error {
	query := ctx.Get("postsQuery").(models.GetPostsQuery)

	slugOrID := ctx.Get("SlugOrID").(string)

	posts, err, msg := threadHandler.Usecase.GetPostsByIdentifier(slugOrID, query)
	if err != nil {
		log.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), *msg)
	}

	response, err := (*posts).MarshalJSON()
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.String(http.StatusOK, string(response))
}

func (threadHandler *ThreadHandler) UpdateThread(ctx echo.Context) error {
	var thread models.Thread
	threadBody := ctx.Get("body").([]byte)
	err := thread.UnmarshalJSON(threadBody)
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	slugOrID := ctx.Get("SlugOrID").(string)

	err, msg := threadHandler.Usecase.UpdateThread(slugOrID, &thread)
	if err != nil {
		log.Error(err)
		//return ctx.JSON(errors.ResolveErrorToCode(err), err.Error())
		return ctx.JSON(errors.ResolveErrorToCode(err), *msg)
	}

	response, err := thread.MarshalJSON()
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(response))
}
