package user

import (
	"github.com/chtvrv/forum_db/app/models"
)

type Usecase interface {
	Create(user *models.User) error
}