package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ReadBody(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		body, err := ioutil.ReadAll(ctx.Request().Body)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		defer ctx.Request().Body.Close()

		ctx.Set("body", body)
		return next(ctx)
	}
}

func ReadUserNickname(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var nickname string
		_, err := fmt.Sscan(ctx.Param("nickname"), &nickname)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		ctx.Set("nickname", nickname)
		return next(ctx)
	}
}

func ReadForumSlug(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var slug string
		_, err := fmt.Sscan(ctx.Param("slug"), &slug)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		ctx.Set("slug", slug)
		return next(ctx)
	}
}

func Headers(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		return next(ctx)
	}
}
