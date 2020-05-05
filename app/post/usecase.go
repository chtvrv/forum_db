package post

import (
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/pkg/errors"
)

type Usecase interface {
	Create(posts *models.Posts, threadIdentifier string) (error, *errors.Message)
}