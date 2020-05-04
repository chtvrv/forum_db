package user

import (
	"github.com/chtvrv/forum_db/app/models"
)

type Repository interface {
	Create(user *models.User) error
	GetUserByNickname(nickname string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	Update(updatedUser *models.User, oldUser *models.User) error
}
