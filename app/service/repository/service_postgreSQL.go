package repository

import (
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/app/service"
	"github.com/jackc/pgx"
	"github.com/labstack/gommon/log"
)

type ServiceStore struct {
	dbConn *pgx.ConnPool
}

func CreateRepository(dbConn_ *pgx.ConnPool) service.Repository {
	return &ServiceStore{dbConn: dbConn_}
}

func (serviceStore *ServiceStore) GetStatus() (*models.Status, error) {
	var status models.Status
	err := serviceStore.dbConn.QueryRow(`SELECT * FROM (SELECT COUNT(*) FROM users) AS usr, (SELECT COUNT(*) FROM forums) AS frm,
		(SELECT COUNT(*) FROM threads) AS thrd, (SELECT COUNT(*) FROM posts) AS pst`).
		Scan(&status.User, &status.Forum, &status.Thread, &status.Post)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &status, nil
}

func (serviceStore *ServiceStore) ClearDB() error {
	res, err := serviceStore.dbConn.Query(`TRUNCATE TABLE users, forums, forum_user, threads, posts, votes CASCADE`)
	defer res.Close()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
