package user

import (
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/pkg/errors"
)

type Usecase interface {
	Create(user *models.User) error
	GetUserByNickname(nickname string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	Update(updatedUser *models.User) (error, *errors.Message)
}
