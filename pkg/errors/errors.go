package errors

import (
	"errors"
	"net/http"
)

type Message struct {
	message string `json:"message"`
}

func CreateMessageNotFoundUser(nickname string) *Message {
	return &Message{
		message: "Can't find user by nickname: " + nickname,
	}
}

func CreateMessageConflictEmail(nickname string) *Message {
	return &Message{
		message: "This email is already registered by user: " + nickname,
	}
}

func CreateMessageNotFoundForum(slug string) *Message {
	return &Message{
		message: "Can't find forum with slug: 1e_gBYzSkiams: " + slug,
	}
}

const (
	Internal     = "Internal error"
	Conflict     = "Conflict with exists data"
	NoPermission = "No permission for current operation"
	NoRows       = "Rows not found"

	UserNotFound = "User not exist"
)

var (
	// общие
	ErrInternal     = errors.New(Internal)
	ErrConflict     = errors.New(Conflict)
	ErrNoPermission = errors.New(NoPermission)
	ErrNoRows       = errors.New(NoRows)

	// ошибки, связанные с пользователем
	ErrUserNotFound = errors.New(UserNotFound)
)

var messToError = map[string]error{
	Internal:     ErrInternal,
	Conflict:     ErrConflict,
	NoPermission: ErrNoPermission,

	UserNotFound: ErrUserNotFound,
	NoRows:       ErrNoRows,
}

var errorToCodeMap = map[error]int{
	// общие
	ErrInternal:     http.StatusInternalServerError,
	ErrConflict:     http.StatusConflict,
	ErrNoPermission: http.StatusForbidden,
	ErrNoRows:       http.StatusNotFound,

	// ошибки, связанные с пользователем
	ErrUserNotFound: http.StatusNotFound,
}

func ResolveErrorToCode(err error) (code int) {
	code, has := errorToCodeMap[err]
	if !has {
		return http.StatusInternalServerError
	}
	return code
}
