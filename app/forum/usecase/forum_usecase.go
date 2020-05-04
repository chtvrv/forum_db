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

func (usecase *ForumUsecase) Create(forum *models.Forum) (error, *errors.Message) {
	user, err := usecase.userRepo.GetUserByNickname(forum.User)
	if err != nil {
		log.Error(err)
		return errors.ErrNoRows, errors.CreateMessageNotFoundUser(forum.User)
	}

	_, err = usecase.forumRepo.GetForumBySlug(forum.Slug)
	if err == nil {
		log.Error("Slug conflict on forum create")
		return errors.ErrConflict, nil
	}
	forum.User = user.Nickname
	err = usecase.forumRepo.Create(forum)
	if err != nil {
		log.Error(err)
		return err, nil
	}
	return nil, nil
}

func (usecase *ForumUsecase) GetForumBySlug(slug string) (*models.Forum, error, *errors.Message) {
	forum, err := usecase.forumRepo.GetForumBySlug(slug)
	if err != nil {
		log.Error(err)
		return nil, err, errors.CreateMessageNotFoundForum(slug)
	}
	return forum, nil, nil
}

func (usecase *ForumUsecase) GetThreadsBySlug(slug string, query models.GetThreadsQuery) (*models.Threads, error, *errors.Message) {
	forum, err, _ := usecase.GetForumBySlug(slug)
	if err != nil {
		log.Error(err)
		return nil, err, errors.CreateMessageNotFoundForum(slug)
	}

	threads, err := usecase.forumRepo.GetThreadsBySlug(forum.Slug, query)
	if err != nil {
		log.Error(err)
		return nil, err, nil
	}

	return threads, nil, nil
}
