package usecase

import (
	"strconv"

	"github.com/chtvrv/forum_db/app/thread"
	"github.com/chtvrv/forum_db/pkg/errors"
	"github.com/labstack/gommon/log"

	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/app/post"
)

type PostUsecase struct {
	postRepo   post.Repository
	threadRepo thread.Repository
}

func CreateUsecase(postRepo_ post.Repository, threadRepo_ thread.Repository) post.Usecase {
	return &PostUsecase{postRepo: postRepo_, threadRepo: threadRepo_}
}

func (usecase *PostUsecase) Create(posts *models.Posts, threadIdentifier string) (error, *errors.Message) {
	var thread *models.Thread
	threadID, err := strconv.Atoi(threadIdentifier)
	if err == nil {
		thread, err = usecase.threadRepo.GetThreadByID(threadID)
	} else {
		thread, err = usecase.threadRepo.GetThreadBySlug(threadIdentifier)
	}
	if thread == nil {
		return errors.ErrNoRows, errors.CreateMessageNotFoundThreadPost(threadID)
	}

	err, msg := usecase.postRepo.Create(posts, thread)
	if err != nil {
		log.Error(err)
		return err, msg
	}
	return nil, nil
}
