package usecase

import (
	"strconv"

	"github.com/chtvrv/forum_db/app/forum"
	"github.com/chtvrv/forum_db/app/thread"
	"github.com/chtvrv/forum_db/app/user"
	"github.com/chtvrv/forum_db/pkg/errors"
	"github.com/labstack/gommon/log"

	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/app/post"
)

type PostUsecase struct {
	postRepo   post.Repository
	threadRepo thread.Repository
	forumRepo  forum.Repository
	userRepo   user.Repository
}

func CreateUsecase(postRepo_ post.Repository, threadRepo_ thread.Repository, forumRepo_ forum.Repository, userRepo_ user.Repository) post.Usecase {
	return &PostUsecase{postRepo: postRepo_, threadRepo: threadRepo_, forumRepo: forumRepo_, userRepo: userRepo_}
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

func (usecase *PostUsecase) GetFullPost(postID int, query models.FullPostQuery) (*models.PostFullInfo, error) {
	var postFullInfo models.PostFullInfo
	var err error

	postFullInfo.Post, err = usecase.postRepo.GetPostByID(postID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	post := postFullInfo.Post

	if query.User {
		postFullInfo.Author, err = usecase.userRepo.GetUserByNickname(post.Author)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if query.Forum {
		postFullInfo.Forum, err = usecase.forumRepo.GetForumBySlug(post.Forum)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if query.Thread {
		postFullInfo.Thread, err = usecase.threadRepo.GetThreadByID(post.Thread)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}

	return &postFullInfo, nil
}

func (usecase *PostUsecase) UpdatePost(updatedPost *models.Post) (error, *errors.Message) {
	var oldPost *models.Post

	oldPost, err := usecase.postRepo.GetPostByID(updatedPost.ID)
	if oldPost == nil {
		return errors.ErrNoRows, errors.CreateMessageNotFoundThreadPost(oldPost.ID)
	}

	err, msg := usecase.postRepo.UpdatePost(updatedPost, oldPost)
	if err != nil {
		log.Error(err)
		return err, msg
	}
	return nil, nil
}
