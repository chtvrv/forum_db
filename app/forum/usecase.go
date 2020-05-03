package forum

import (
	"github.com/chtvrv/forum_db/app/models"
)

type Usecase interface {
	Create(forum *models.Forum) error
	GetForumBySlug(slug string) (*models.Forum, error)
}
