package delivery

import (
	"net/http"

	"github.com/chtvrv/forum_db/app/middleware"
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/pkg/errors"

	"github.com/chtvrv/forum_db/app/forum"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type ForumHandler struct {
	Usecase forum.Usecase
}

func CreateHandler(router *echo.Echo, usecase forum.Usecase) {
	handler := &ForumHandler{
		Usecase: usecase,
	}

	router.POST("/api/forum/create", handler.Create, middleware.ReadBody, middleware.Headers)
	router.GET("/api/forum/:slug/details", handler.Get, middleware.ReadForumSlug, middleware.Headers)
	router.GET("/api/forum/:slug/threads", handler.GetThreads,
		middleware.ReadForumSlug, middleware.ReadGetThreadsQuery, middleware.Headers)
	router.GET("/api/forum/:slug/users", handler.GetUsers,
		middleware.ReadForumSlug, middleware.ReadGetThreadsQuery, middleware.Headers)
}

func (forumHandler *ForumHandler) Create(ctx echo.Context) error {
	var forum models.Forum
	forumBody := ctx.Get("body").([]byte)
	err := forum.UnmarshalJSON(forumBody)
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	err, msg := forumHandler.Usecase.Create(&forum)
	// Успешно создали
	if err == nil {
		response, err := forum.MarshalJSON()
		if err != nil {
			log.Error(err)
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.String(http.StatusCreated, string(response))
	}
	// Конфликт
	if err == errors.ErrConflict {
		previousForum, _, _ := forumHandler.Usecase.GetForumBySlug(forum.Slug)
		response, err := previousForum.MarshalJSON()
		if err != nil {
			log.Error(err)
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.String(http.StatusConflict, string(response))
	}
	//Нет пользователя
	if err == errors.ErrNoRows {
		return ctx.JSON(errors.ResolveErrorToCode(err), *msg)
	}
	// Незарегистрированная ошибка
	log.Error(err)
	return ctx.String(errors.ResolveErrorToCode(err), err.Error())
}

func (forumHandler *ForumHandler) Get(ctx echo.Context) error {
	slug := ctx.Get("slug").(string)
	forum, err, msg := forumHandler.Usecase.GetForumBySlug(slug)
	if err != nil {
		log.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), *msg)
	}

	response, err := forum.MarshalJSON()
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(response))
}

func (forumHandler *ForumHandler) GetThreads(ctx echo.Context) error {
	query := ctx.Get("threadsQuery").(models.GetThreadsQuery)

	forumSlug := ctx.Get("slug").(string)
	threads, err, msg := forumHandler.Usecase.GetThreadsBySlug(forumSlug, query)
	if err != nil {
		log.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), *msg)
	}

	response, err := (*threads).MarshalJSON()
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.String(http.StatusOK, string(response))
}

func (forumHandler *ForumHandler) GetUsers(ctx echo.Context) error {
	query := ctx.Get("threadsQuery").(models.GetThreadsQuery)

	forumSlug := ctx.Get("slug").(string)
	users, err, msg := forumHandler.Usecase.GetUsersBySlug(forumSlug, query)
	if err != nil {
		log.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), *msg)
	}

	response, err := (*users).MarshalJSON()
	if err != nil {
		log.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.String(http.StatusOK, string(response))
}
