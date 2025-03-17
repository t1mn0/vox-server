package test_storage

import (
	"vox-server/internal/storage"
)

type InMemoryStorage struct {
	userRepository *UserRepository
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{}
}

func (storage *InMemoryStorage) Users() storage.UserRepository {
	if storage.userRepository != nil {
		return storage.userRepository
	}

	storage.userRepository = NewUserRepository()

	return *storage.userRepository
}
