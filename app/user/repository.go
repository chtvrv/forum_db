package user

import (
	"github.com/chtvrv/forum_db/app/models"
)

type Repository interface {
	Create(user *models.User) error
}
