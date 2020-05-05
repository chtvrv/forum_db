package repository

import (
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/app/thread"
	"github.com/chtvrv/forum_db/pkg/errors"
	"github.com/jackc/pgx"
	"github.com/labstack/gommon/log"
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

	result, err := threadStore.dbConn.Exec(`UPDATE forums SET threads = threads + $1 WHERE slug = $2`, 1, thread.Forum)
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
	err := threadStore.dbConn.QueryRow(`SELECT * FROM threads WHERE id = $1`, id).
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
