package storage

import "git-server/internal/models"

type UserRepository interface {
	Create(*models.User) error
	FindByLogin(string) (*models.User, error)
}
