package teststorage

import (
	"errors"
	"git-server/internal/models"
)

type UserRepository struct {
	storage *Storage
	users   map[string]*models.User
}

func (userrepo *UserRepository) Create(user *models.User) error {
	if err := user.Validate(true); err != nil {
		return err
	}

	if err := user.BeforeCreate(); err != nil {
		return err
	}

	userrepo.users[user.Login] = user

	return nil
}

func (userrepo *UserRepository) FindByLogin(login string) (*models.User, error) {
	user, ok := userrepo.users[login]

	if ok {
		return user, nil
	}

	return nil, errors.New("user not found")
}
