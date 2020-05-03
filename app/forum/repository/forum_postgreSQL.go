package repository

import (
	"github.com/chtvrv/forum_db/app/forum"
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/pkg/errors"
	"github.com/jackc/pgx"
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
