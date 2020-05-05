package forum

import (
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/pkg/errors"
)

type Usecase interface {
	Create(forum *models.Forum) (error, *errors.Message)
	GetForumBySlug(slug string) (*models.Forum, error, *errors.Message)
	GetThreadsBySlug(slug string, query models.GetThreadsQuery) (*models.Threads, error, *errors.Message)
	GetUsersBySlug(slug string, query models.GetThreadsQuery) (*models.Users, error, *errors.Message)
}
