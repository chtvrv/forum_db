package usecase

import (
	"github.com/chtvrv/forum_db/app/forum"
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/app/thread"
	"github.com/chtvrv/forum_db/app/user"
	"github.com/chtvrv/forum_db/pkg/errors"
	"github.com/labstack/gommon/log"
)

type ThreadUsecase struct {
	threadRepo thread.Repository
	userRepo   user.Repository
	forumRepo  forum.Repository
}

func CreateUsecase(threadRepo_ thread.Repository, userRepo_ user.Repository, forumRepo_ forum.Repository) thread.Usecase {
	return &ThreadUsecase{threadRepo: threadRepo_, userRepo: userRepo_, forumRepo: forumRepo_}
}

func (usecase *ThreadUsecase) Create(thread *models.Thread) (error, *errors.Message) {
	_, err := usecase.userRepo.GetUserByNickname(thread.Author)
	if err != nil {
		log.Error(err)
		return err, errors.CreateMessageNotFoundThreadAuthor(thread.Author)
	}

	prevForum, err := usecase.forumRepo.GetForumBySlug(thread.Forum)
	if err != nil {
		log.Error(err)
		return err, errors.CreateMessageNotFoundThreadForum(thread.Forum)
	}

	if thread.Slug != "" {
		oldThread, _ := usecase.threadRepo.GetThreadBySlug(thread.Slug)
		if oldThread != nil {
			log.Error("Thread slug conflict on thread create")
			return errors.ErrConflict, nil
		}
	}
	thread.Forum = prevForum.Slug
	err = usecase.threadRepo.Create(thread)
	if err != nil {
		log.Error(err)
		return err, nil
	}
	return nil, nil
}

func (usecase *ThreadUsecase) GetThreadBySlug(slug string) (*models.Thread, error) {
	thread, err := usecase.threadRepo.GetThreadBySlug(slug)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return thread, nil
}
