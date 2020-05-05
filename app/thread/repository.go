package thread

import (
	"github.com/chtvrv/forum_db/app/models"
)

type Repository interface {
	Create(thread *models.Thread) error
	GetThreadBySlug(slug string) (*models.Thread, error)
	GetThreadByID(id int) (*models.Thread, error)
}
