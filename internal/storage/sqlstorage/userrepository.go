// light copypast from sqlc

package sqlstorage

import (
	"git-server/internal/models"
	_ "git-server/internal/storage"
)

type UserRepo struct {
	storage *DBStorage
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (login, username, email, encrypted_password)
VALUES ($1, $2, $3, $4)
RETURNING login`

func (userrepo UserRepo) Create(arg_user *models.User) error {
	if err := arg_user.Validate(true); err != nil {
		return err
	}

	if err := arg_user.BeforeCreate(); err != nil {
		return err
	}

	row := userrepo.storage.db.QueryRow(
		createUser,
		arg_user.Login,
		arg_user.Username,
		arg_user.Email,
		arg_user.EncryptedPassword,
	)

	err := row.Scan(
		&arg_user.Login,
	)

	if err != nil {
		return err
	}

	return nil
}

const findUserByLogin = `-- name: FindByLogin :one
SELECT login, username, email, encrypted_password FROM users
WHERE login = $1`

func (userrepo UserRepo) FindByLogin(login string) (*models.User, error) {
	row := userrepo.storage.db.QueryRow(findUserByLogin, login)
	var u models.User
	err := row.Scan(
		&u.Login,
		&u.Username,
		&u.Email,
		&u.EncryptedPassword,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// other sql queries...
