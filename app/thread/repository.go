package thread

import (
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/pkg/errors"
)

type Repository interface {
	Create(thread *models.Thread) error
	GetThreadBySlug(slug string) (*models.Thread, error)
	GetThreadByID(id int) (*models.Thread, error)
	VoteForThread(vote *models.Vote) (error, *errors.Message)
	GetPostsByThread(thread *models.Thread, query models.GetPostsQuery) (*models.Posts, error, *errors.Message)
	UpdateThread(updatedThread *models.Thread, oldThread *models.Thread) (error, *errors.Message)
}
