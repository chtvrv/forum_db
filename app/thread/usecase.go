package thread

import (
	"github.com/chtvrv/forum_db/app/models"
)

type Usecase interface {
	Create(thread *models.Thread) error
}