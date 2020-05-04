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

func (userStore *UserStore) Update(updatedUser *models.User, oldUser *models.User) error {
	if updatedUser.Fullname == "" && updatedUser.About == "" && updatedUser.Email == "" {
		*updatedUser = *oldUser
		return nil
	}

	if updatedUser.Fullname != "" {
		oldUser.Fullname = updatedUser.Fullname
	}

	if updatedUser.Email != "" {
		oldUser.Email = updatedUser.Email
	}

	if updatedUser.About != "" {
		oldUser.About = updatedUser.About
	}

	_, err := userStore.dbConn.Exec("UPDATE users SET fullname = $1, email = $2, about = $3 WHERE nickname = $4",
		oldUser.Fullname, oldUser.Email, oldUser.About, oldUser.Nickname)
	if err != nil {
		log.Error(err)
		return err
	}

	*updatedUser = *oldUser
	return nil
}
