package storage

import "vox-server/internal/models"

type UserRepository interface {
	Count() int
	IsEmpty() bool
	Create(*models.User) error
	FindByLogin(string) (*models.User, error)
	FindByEmail(string) (*models.User, error)
	DeleteByLogin(login string) error
	DeleteByEmail(email string) error
	Update(user *models.User) error
}
