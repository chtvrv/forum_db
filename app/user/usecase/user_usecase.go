package usecase

import (
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/app/user"
	"github.com/chtvrv/forum_db/pkg/errors"
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

func (usecase *UserUsecase) Update(updatedUser *models.User) (error, *errors.Message) {
	oldUser, _ := usecase.userRepo.GetUserByNickname(updatedUser.Nickname)
	if oldUser == nil {
		return errors.ErrUserNotFound, errors.CreateMessageNotFoundUser(updatedUser.Nickname)
	}

	userWithThisEmail, _ := usecase.userRepo.GetUserByEmail(updatedUser.Email)
	if userWithThisEmail != nil && userWithThisEmail.Nickname != oldUser.Nickname {
		return errors.ErrConflict, errors.CreateMessageConflictEmail(userWithThisEmail.Nickname)
	}

	err := usecase.userRepo.Update(updatedUser, oldUser)
	if err != nil {
		log.Error(err)
		return err, nil
	}
	return nil, nil
}
