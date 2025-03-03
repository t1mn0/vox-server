package teststorage

import (
	"git-server/internal/models"
	"git-server/internal/storage"
)

type Storage struct {
	userRepository *UserRepository
}

func New() *Storage {
	return &Storage{}
}

func (st *Storage) User() storage.UserRepository {
	if st.userRepository != nil {
		return st.userRepository
	}

	st.userRepository = &UserRepository{
		storage: st,
		users:   make(map[string]*models.User),
	}

	return st.userRepository
}
