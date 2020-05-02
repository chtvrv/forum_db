package usecase

import (
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/app/user"
	"github.com/labstack/gommon/log"
)

type UserUsecase struct {
	userRepo user.Repository
}

func CreateUsecase(userRepo_ user.Repository) user.Usecase {
	return &UserUsecase{userRepo: userRepo_}
}

func (usecase *UserUsecase) Create(user *models.User) error {
	err := usecase.userRepo.Create(user)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
