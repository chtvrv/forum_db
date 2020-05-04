package repository

import (
	"fmt"

	"github.com/chtvrv/forum_db/app/forum"
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/pkg/errors"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/labstack/gommon/log"
)

type ForumStore struct {
	dbConn *pgx.ConnPool
}

func CreateRepository(dbConn_ *pgx.ConnPool) forum.Repository {
	return &ForumStore{dbConn: dbConn_}
}

func (forumStore *ForumStore) Create(forum *models.Forum) error {
	result, err := forumStore.dbConn.Exec(`INSERT INTO forums (title, usr, slug) VALUES ($1, $2, $3)`,
		forum.Title, forum.User, forum.Slug)

	if err != nil {
		log.Error(err)
		return errors.ErrConflict
	}

	if result.RowsAffected() != 1 {
		log.Error("Forum data collision on create")
		return errors.ErrConflict
	}

	result, err = forumStore.dbConn.Exec(`INSERT INTO forum_user (slug, nickname) VALUES ($1, $2)`,
		forum.Slug, forum.User)

	if err != nil {
		log.Error(err)
		return errors.ErrConflict
	}

	if result.RowsAffected() != 1 {
		log.Error("Forum data collision on create")
		return errors.ErrConflict
	}

	return nil
}

func (forumStore *ForumStore) GetForumBySlug(slug string) (*models.Forum, error) {
	var forum models.Forum
	err := forumStore.dbConn.QueryRow(`SELECT * FROM forums WHERE slug = $1`, slug).
		Scan(&forum.Slug, &forum.Title, &forum.User)

	if err != nil && err != pgx.ErrNoRows {
		log.Error(err)
		return nil, err
	}

	if err == pgx.ErrNoRows {
		return nil, errors.ErrNoRows
	}

	return &forum, nil
}

func (forumStore *ForumStore) GetThreadsBySlug(slug string, query models.GetThreadsQuery) (*models.Threads, error) {
	dbQuery := CreateThreadsQuery(query)
	result, err := forumStore.dbConn.Query(dbQuery, slug)
	if err != nil {
		log.Error(err)
		return nil, errors.ErrInternal
	}

	defer result.Close()

	var threads models.Threads
	for result.Next() {
		var thread models.Thread
		slug := &pgtype.Varchar{}

		err := result.Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, slug, &thread.Created)
		if err != nil {
			log.Error(err)
			return nil, errors.ErrInternal

		}

		thread.Slug = slug.String
		threads = append(threads, thread)
	}

	return &threads, nil
}

func CreateThreadsQuery(threadsQuery models.GetThreadsQuery) string {
	sinceToken := ``
	sortToken := ` ORDER BY created`
	limitToken := fmt.Sprintf(` LIMIT %d`, threadsQuery.Limit)

	if threadsQuery.Since != "" {
		sinceToken = ` AND created `
		if threadsQuery.Desc {
			sinceToken += `<= `
		} else {
			sinceToken += `>= `
		}
		sinceToken += `'` + threadsQuery.Since + `'`
	}

	if threadsQuery.Desc {
		sortToken += ` DESC`
	}

	return `SELECT * FROM threads WHERE forum = $1` + sinceToken + sortToken + limitToken
}
