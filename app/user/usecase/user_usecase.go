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

func (usecase *UserUsecase) GetUserByNickname(nickname string) (*models.User, error) {
	user, err := usecase.userRepo.GetUserByNickname(nickname)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return user, nil
}

func (usecase *UserUsecase) GetUserByEmail(email string) (*models.User, error) {
	user, err := usecase.userRepo.GetUserByEmail(email)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return user, nil
}
