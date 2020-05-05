package repository

import (
	"time"

	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/app/post"
	"github.com/chtvrv/forum_db/pkg/errors"
	"github.com/jackc/pgx"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
)

type PostStore struct {
	dbConn *pgx.ConnPool
}

func CreateRepository(dbConn_ *pgx.ConnPool) post.Repository {
	return &PostStore{dbConn: dbConn_}
}

func (postStore *PostStore) Create(posts *models.Posts, thread *models.Thread) (error, *errors.Message) {
	postParents, msg := postStore.CheckAndGetParents(posts, thread)
	if msg != nil {
		return errors.ErrConflict, msg
	}

	generatedIDs := postStore.GeneratePostIDSequence(len(*posts))
	if generatedIDs == nil {
		return errors.ErrInternal, nil
	}

	currentTimestamp, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Error(err)
		return errors.ErrInternal, nil
	}

	// Сама вставка
	for id, post := range *posts {
		post.ID = (*generatedIDs)[id]
		post.Forum = thread.Forum
		post.Thread = thread.ID
		post.Created = currentTimestamp
		post.Path = append((*postParents)[post.Parent].Path, int64((*generatedIDs)[id]))
		(*posts)[id] = post

		result, err := postStore.dbConn.Exec(`INSERT INTO posts (id, parent, author, message, forum, thread, created, path) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			post.ID, post.Parent, post.Author, post.Message, post.Forum, post.Thread, post.Created, pq.Array(post.Path))
		if err != nil || result.RowsAffected() == 0 {
			log.Error(err)
			return errors.ErrNoRows, errors.CreateNotFoundAuthorPost(post.Author)
		}

		_, _ = postStore.dbConn.Exec(`INSERT INTO forum_user (slug, nickname) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			thread.Forum, post.Author)
	}

	result, err := postStore.dbConn.Exec(`UPDATE forums SET posts = posts + $1 WHERE slug = $2`, len(*posts), thread.Forum)
	if err != nil || result.RowsAffected() == 0 {
		return errors.ErrInternal, nil
	}

	return nil, nil
}

func (postStore *PostStore) GetPostByID(id int) (*models.Post, error) {
	var post models.Post
	err := postStore.dbConn.QueryRow(`SELECT * FROM posts WHERE id = $1`, id).
		Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited,
			&post.Forum, &post.Thread, &post.Created, pq.Array(&post.Path))

	if err != nil && err != pgx.ErrNoRows {
		log.Error(err)
		return nil, err
	}

	if err == pgx.ErrNoRows {
		return nil, errors.ErrNoRows
	}

	return &post, nil
}

func (postStore *PostStore) CheckAndGetParents(posts *models.Posts, thread *models.Thread) (*map[int]models.Post, *errors.Message) {
	postParents := make(map[int]models.Post)

	for _, post := range *posts {
		if _, exist := postParents[post.Parent]; !exist && post.Parent != 0 {
			parent, _ := postStore.GetPostByID(post.Parent)
			if parent == nil || parent.Thread != thread.ID {
				return nil, errors.CreateMessageConflictCreatePost()
			}
			postParents[post.Parent] = *parent
		}
	}

	return &postParents, nil
}

func (postStore *PostStore) GeneratePostIDSequence(sequenceLength int) *[]int {
	sequenceRows, err := postStore.dbConn.Query(`SELECT nextval(pg_get_serial_sequence('posts', 'id')) FROM generate_series(1, $1)`, sequenceLength)
	if err != nil {
		return nil
	}

	defer sequenceRows.Close()
	generatedIDs := make([]int, 0)
	for sequenceRows.Next() {
		var id int
		err = sequenceRows.Scan(&id)
		if err != nil {
			return nil
		}
		generatedIDs = append(generatedIDs, id)
	}

	return &generatedIDs
}

func (postStore *PostStore) UpdatePost(updatedPost *models.Post, oldPost *models.Post) (error, *errors.Message) {
	if updatedPost.Message == "" || updatedPost.Message == oldPost.Message {
		*updatedPost = *oldPost
		return nil, nil
	}

	oldPost.Message = updatedPost.Message
	oldPost.IsEdited = true

	_, err := postStore.dbConn.Exec("UPDATE posts SET message = $1, is_edited = $2 WHERE id = $3",
		oldPost.Message, oldPost.IsEdited, oldPost.ID)
	if err != nil {
		log.Error(err)
		return err, nil
	}

	*updatedPost = *oldPost
	return nil, nil
}
