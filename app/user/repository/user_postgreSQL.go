package repository

import (
	"github.com/chtvrv/forum_db/app/models"
	"github.com/chtvrv/forum_db/app/user"
	"github.com/chtvrv/forum_db/pkg/errors"
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
		return errors.ErrConflict
	}

	if result.RowsAffected() != 1 {
		log.Error("User data collision on create")
		return errors.ErrConflict
	}

	return nil
}

func (userStore *UserStore) GetUserByNickname(nickname string) (*models.User, error) {
	var user models.User
	err := userStore.dbConn.QueryRow(`SELECT * FROM users WHERE nickname = $1`, nickname).
		Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)

	if err != nil && err != pgx.ErrNoRows {
		log.Error(err)
		return nil, err
	}

	if err == pgx.ErrNoRows {
		return nil, errors.ErrNoRows
	}

	return &user, nil
}

func (userStore *UserStore) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := userStore.dbConn.QueryRow(`SELECT * FROM users WHERE email = $1`, email).
		Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)

	if err != nil && err != pgx.ErrNoRows {
		log.Error(err)
		return nil, err
	}

	if err == pgx.ErrNoRows {
		return nil, errors.ErrNoRows
	}

	return &user, nil
}
