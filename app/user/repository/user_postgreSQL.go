package repository

import (
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/app/user"
	"github.com/jackc/pgx"
	"github.com/labstack/gommon/log"
)

type UserStore struct {
	dbConn *pgx.ConnPool
}

func CreateRepository(dbConn_ *pgx.ConnPool) user.Repository {
	return &UserStore{dbConn: dbConn_}
}

func (userStore *UserStore) Create(user *models.User) error {
	result, err := userStore.dbConn.Exec(`INSERT INTO users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4)`,
		user.Nickname, user.Fullname, user.Email, user.About)

	if err != nil {
		log.Error(err)
		// TODO модель ошибок
		return nil
	}

	if result.RowsAffected() != 1 {
		log.Error("Collision!")
		return nil
	}

	return nil
}
