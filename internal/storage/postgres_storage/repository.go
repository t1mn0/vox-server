package postgres_storage

import (
	"vox-server/internal/models"
)

type UserRepository struct {
	storage *DBStorage
}

const countUsers = `-- name: CountUsers :one
SELECT COUNT(*) FROM users`

func (repository UserRepository) Count() int {
	var count int
	err := repository.storage.db.QueryRow(countUsers).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func (repository UserRepository) IsEmpty() bool {
	return repository.Count() == 0
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (login, username, email, encrypted_password)
VALUES ($1, $2, $3, $4)
RETURNING login`

func (repository UserRepository) Create(arg_user *models.User) error {
	if err := arg_user.Validate(true); err != nil {
		return err
	}

	if err := arg_user.BeforeCreate(); err != nil {
		return err
	}

	row := repository.storage.db.QueryRow(
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

func (repository UserRepository) FindByLogin(login string) (*models.User, error) {
	row := repository.storage.db.QueryRow(findUserByLogin, login)
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

const findUserByEmail = `-- name: FindByEmail :one
SELECT login, username, email, encrypted_password FROM users
WHERE email = $1`

func (repository UserRepository) FindByEmail(email string) (*models.User, error) {
	row := repository.storage.db.QueryRow(findUserByEmail, email)
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

const deleteUserByLogin = `-- name: DeleteByLogin :exec
DELETE FROM users WHERE login = $1`

func (repository UserRepository) DeleteByLogin(login string) error {
	_, err := repository.storage.db.Exec(deleteUserByLogin, login)
	return err
}

const deleteUserByEmail = `-- name: DeleteByEmail :exec
DELETE FROM users WHERE email = $1`

func (repository UserRepository) DeleteByEmail(email string) error {
	_, err := repository.storage.db.Exec(deleteUserByEmail, email)
	return err
}

// const updateUser = `-- name: UpdateUser :exec
// UPDATE users
// SET username = $2, email = $3, encrypted_password = $4
// WHERE login = $1`

func (repository UserRepository) Update(user *models.User) error {
	panic("unimplemented!")
	// user.Validate(bool); err != nil { . . . }
	// if it was password => user.BeforeCreate(); err != nil { . . . }
}
