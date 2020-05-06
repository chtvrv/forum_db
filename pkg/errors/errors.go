package errors

import (
	"errors"
	"net/http"
	"strconv"
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
		message: "Can't find forum with slug: " + slug,
	}
}
func CreateMessageConflictCreatePost() *Message {
	return &Message{
		message: "Parent post was created in another thread",
	}
}

func CreateMessageNotFoundThreadPost(id int) *Message {
	return &Message{
		message: "Can't find post thread by id: " + strconv.Itoa(id),
	}
}

func CreateMessageNotFoundThreadAuthor(nickname string) *Message {
	return &Message{
		message: "Can't find thread author by nickname: " + nickname,
	}
}

func CreateMessageNotFoundThreadForum(slug string) *Message {
	return &Message{
		message: "Can't find thread forum by slug: " + slug,
	}
}

func CreateNotFoundAuthorPost(nickname string) *Message {
	return &Message{
		message: "Can't find post author by nickname: " + nickname,
	}
}

func CreateNotFoundPost(id string) *Message {
	return &Message{
		message: "Can't find post with id: " + id,
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
