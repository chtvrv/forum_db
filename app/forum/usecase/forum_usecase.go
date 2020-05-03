package usecase

import (
	"github.com/chtvrv/forum_db/app/forum"
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/app/user"
	"github.com/chtvrv/forum_db/pkg/errors"
	"github.com/labstack/gommon/log"
)

type ForumUsecase struct {
	forumRepo forum.Repository
	userRepo  user.Repository
}

func CreateUsecase(forumRepo_ forum.Repository, userRepo_ user.Repository) forum.Usecase {
	return &ForumUsecase{forumRepo: forumRepo_, userRepo: userRepo_}
}

func (usecase *ForumUsecase) Create(forum *models.Forum) error {
	_, err := usecase.userRepo.GetUserByNickname(forum.User)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = usecase.forumRepo.GetForumBySlug(forum.Slug)
	if err == nil {
		log.Error("Slug conflict on forum create")
		return errors.ErrConflict
	}

	err = usecase.forumRepo.Create(forum)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (usecase *ForumUsecase) GetForumBySlug(slug string) (*models.Forum, error) {
	forum, err := usecase.forumRepo.GetForumBySlug(slug)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return forum, nil
}
