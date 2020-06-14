package repository

import (
	"fmt"

	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/app/thread"
	"github.com/chtvrv/forum_db/pkg/errors"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
)

type ThreadStore struct {
	dbConn *pgx.ConnPool
}

func CreateRepository(dbConn_ *pgx.ConnPool) thread.Repository {
	return &ThreadStore{dbConn: dbConn_}
}

func (threadStore *ThreadStore) Create(thread *models.Thread) error {
	// slug трэда уникален, но при этом опционален
	var err error
	if thread.Slug == "" {
		err = threadStore.dbConn.QueryRow(`INSERT INTO threads (title, author, forum, message, created) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			thread.Title, thread.Author, thread.Forum, thread.Message, thread.Created).Scan(&thread.ID)
	} else {
		err = threadStore.dbConn.QueryRow(`INSERT INTO threads (title, author, forum, message, created, slug) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			thread.Title, thread.Author, thread.Forum, thread.Message, thread.Created, thread.Slug).Scan(&thread.ID)
	}

	result, err := threadStore.dbConn.Exec(`INSERT INTO forum_user (slug, nickname) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		thread.Forum, thread.Author)

	if err != nil {
		log.Error(err)
		return errors.ErrConflict
	}

	result, err = threadStore.dbConn.Exec(`UPDATE forums SET threads = threads + $1 WHERE slug = $2`, 1, thread.Forum)
	if err != nil || result.RowsAffected() == 0 {
		return errors.ErrInternal
	}

	if err != nil {
		log.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (threadStore *ThreadStore) GetThreadBySlug(slug string) (*models.Thread, error) {
	var thread models.Thread
	err := threadStore.dbConn.QueryRow(`SELECT * FROM threads WHERE slug = $1`, slug).
		Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)

	if err != nil && err != pgx.ErrNoRows {
		log.Error(err)
		return nil, err
	}

	if err == pgx.ErrNoRows {
		return nil, errors.ErrNoRows
	}

	return &thread, nil
}

func (threadStore *ThreadStore) GetThreadByID(id int) (*models.Thread, error) {
	var thread models.Thread
	nullSlug := &pgtype.Varchar{}

	err := threadStore.dbConn.QueryRow(`SELECT * FROM threads WHERE id = $1`, id).
		Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, nullSlug, &thread.Created)

	if err != nil && err != pgx.ErrNoRows {
		log.Error(err)
		return nil, err
	}

	if err == pgx.ErrNoRows {
		return nil, errors.ErrNoRows
	}

	thread.Slug = nullSlug.String
	return &thread, nil
}

func (threadStore *ThreadStore) VoteForThread(vote *models.Vote) (error, *errors.Message) {
	var oldVote models.Vote
	result, err := threadStore.dbConn.Query(`SELECT voice FROM votes WHERE nickname = $1 AND thread = $2`, vote.Nickname, vote.Thread)
	defer result.Close()

	// Голос уже существует, обновляем его
	if err == nil && result.Next() {
		result.Scan(&oldVote.Voice)
		if oldVote.Voice == vote.Voice {
			return nil, nil
		}

		_, err := threadStore.dbConn.Exec(`UPDATE votes SET voice = $1 WHERE nickname = $2 AND thread = $3`, vote.Voice, vote.Nickname, vote.Thread)
		if err != nil {
			log.Error(err)
			return errors.ErrInternal, nil
		}

		// Для пересчёта разницы при обновлении threads таблицы
		vote.Voice *= 2

		// Голосуем впервые
	} else {
		_, err := threadStore.dbConn.Exec(`INSERT INTO votes (nickname, voice, thread) VALUES ($1, $2, $3)`, vote.Nickname, vote.Voice, vote.Thread)
		if err != nil {
			log.Error(err)
			return errors.ErrInternal, nil
		}
	}

	_, err = threadStore.dbConn.Exec(`UPDATE threads SET votes = votes + $1 WHERE id = $2`, vote.Voice, vote.Thread)
	if err != nil {
		log.Error(err)
		return errors.ErrInternal, nil
	}

	return nil, nil
}

func (threadStore *ThreadStore) GetPostsByThread(thread *models.Thread, query models.GetPostsQuery) (*models.Posts, error, *errors.Message) {
	var dbQuery string
	if query.Sort == "flat" {
		dbQuery = CreatePostsFlatQuery(query)
	} else if query.Sort == "tree" {
		dbQuery = CreatePostsTreeQuery(query)
	} else {
		dbQuery = CreatePostsParentTreeQuery(query)
	}

	result, err := threadStore.dbConn.Query(dbQuery, thread.ID)
	defer result.Close()
	if err != nil {
		log.Error(err)
		return nil, err, nil
	}

	posts := make(models.Posts, 0)
	for result.Next() {
		var post models.Post

		err := result.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum,
			&post.Thread, &post.Created, pq.Array(&post.Path))
		if err != nil {
			log.Error(err)
			return nil, errors.ErrInternal, nil
		}
		posts = append(posts, post)
	}

	return &posts, nil, nil
}

func CreatePostsFlatQuery(postsQuery models.GetPostsQuery) string {
	sinceToken := ``
	sortToken := ` ORDER BY id`
	limitToken := fmt.Sprintf(` LIMIT %d`, postsQuery.Limit)

	if postsQuery.Since != "" {
		sinceToken = ` AND id `
		if postsQuery.Desc {
			sinceToken += `< `
		} else {
			sinceToken += `> `
		}
		sinceToken += `'` + postsQuery.Since + `'`
	}

	if postsQuery.Desc {
		sortToken += ` DESC`
	}

	return `SELECT * FROM posts WHERE thread = $1` + sinceToken + sortToken + limitToken
}

func CreatePostsTreeQuery(postsQuery models.GetPostsQuery) string {
	sinceToken := ``
	sortToken := ` ORDER BY`
	limitToken := fmt.Sprintf(` LIMIT %d`, postsQuery.Limit)

	if postsQuery.Since != "" {
		sinceToken = ` AND path `
		if postsQuery.Desc {
			sinceToken += `< `
		} else {
			sinceToken += `> `
		}
		sinceToken += "(SELECT path FROM posts WHERE id = " + postsQuery.Since + `)`
	}

	if postsQuery.Desc {
		sortToken += ` path DESC, id DESC`
	} else {
		sortToken += ` path, id`
	}

	return `SELECT * FROM posts WHERE thread = $1` + sinceToken + sortToken + limitToken
}

func CreatePostsParentTreeQuery(postsQuery models.GetPostsQuery) string {
	sinceToken := ``
	sortTokenInner := ` ORDER BY id`
	sortTokenExternal := ` ORDER BY`
	limitToken := fmt.Sprintf(` LIMIT %d`, postsQuery.Limit)

	if postsQuery.Since != "" {
		sinceToken = ` AND path[1] `
		if postsQuery.Desc {
			sinceToken += `< `
		} else {
			sinceToken += `> `
		}
		sinceToken += "(SELECT path[1] FROM posts WHERE id = " + postsQuery.Since + `)`
	}

	if postsQuery.Desc {
		sortTokenInner += ` DESC`
		sortTokenExternal += ` path[1] DESC, path, id`
	} else {
		sortTokenExternal += ` path`
	}

	return `SELECT * FROM posts WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0` + sinceToken +
		sortTokenInner + limitToken + ")" + sortTokenExternal
}

func (threadStore *ThreadStore) UpdateThread(updatedThread *models.Thread, oldThread *models.Thread) (error, *errors.Message) {
	if updatedThread.Title == "" && updatedThread.Message == "" {
		*updatedThread = *oldThread
		return nil, nil
	}

	if updatedThread.Title != "" {
		oldThread.Title = updatedThread.Title
	}

	if updatedThread.Message != "" {
		oldThread.Message = updatedThread.Message
	}

	_, err := threadStore.dbConn.Exec("UPDATE threads SET title = $1, message = $2 WHERE id = $3",
		oldThread.Title, oldThread.Message, oldThread.ID)
	if err != nil {
		log.Error(err)
		return err, nil
	}

	*updatedThread = *oldThread
	return nil, nil
}
