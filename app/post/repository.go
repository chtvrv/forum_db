package post

import (
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/pkg/errors"
)

type Repository interface {
	Create(posts *models.Posts, thread *models.Thread) (error, *errors.Message)
	GetPostByID(id int) (*models.Post, error)
	UpdatePost(updatedPost *models.Post, oldPost *models.Post) (error, *errors.Message)
}
