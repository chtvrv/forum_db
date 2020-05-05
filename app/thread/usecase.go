package thread

import (
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/pkg/errors"
)

type Usecase interface {
	Create(thread *models.Thread) (error, *errors.Message)
	GetThreadByIdentifier(identifier string) (*models.Thread, error, *errors.Message)
	GetThreadBySlug(slug string) (*models.Thread, error)
	VoteForThread(vote *models.Vote, threadIdentifier string) (*models.Thread, error, *errors.Message)
	GetPostsByIdentifier(threadIdentifier string, query models.GetPostsQuery) (*models.Posts, error, *errors.Message)
	UpdateThread(threadIdentifier string, updatedThread *models.Thread) (error, *errors.Message)
}
