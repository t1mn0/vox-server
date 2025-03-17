package test_storage

import (
	"fmt"
	"sync"
	"vox-server/internal/models"
)

type UserRepository struct {
	users  map[string]*models.User // login -> user
	emails map[string]string       // email -> login
	mu     *sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:  make(map[string]*models.User),
		emails: make(map[string]string),
		mu:     &sync.RWMutex{},
	}
}

func (repository UserRepository) Count() int {
	return len(repository.users)
}

func (repository UserRepository) IsEmpty() bool {
	return len(repository.users) == 0
}

func (repository UserRepository) Create(user *models.User) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	if _, ok := repository.users[user.Login]; ok {
		return fmt.Errorf("user with such login '%s' already exist", user.Login)
	}

	if _, ok := repository.emails[user.Email]; ok {
		return fmt.Errorf("user with such email '%s' already exist", user.Email)
	}

	if err := user.Validate(true); err != nil {
		return err
	}

	if err := user.BeforeCreate(); err != nil {
		return err
	}

	repository.users[user.Login] = user
	repository.emails[user.Email] = user.Login

	return nil
}

func (repository UserRepository) FindByLogin(login string) (*models.User, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	user, ok := repository.users[login]

	if ok {
		return user, nil
	}

	return nil, fmt.Errorf("user with login '%s' not found", login)
}

func (repository UserRepository) FindByEmail(email string) (*models.User, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	user, ok := repository.users[repository.emails[email]]

	if ok {
		return user, nil
	}

	return nil, fmt.Errorf("user with email '%s' not found", email)
}

// O(n) search pair (email -> login) in repository.emails
func (repository UserRepository) DeleteByLogin(login string) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	user, ok := repository.users[login]
	if !ok {
		return fmt.Errorf("user with login '%s' not found", login)
	}

	delete(repository.users, login)

	email := user.Email
	for e, l := range repository.emails {
		if l == login {
			email = e
			break
		}
	}
	delete(repository.emails, email)

	return nil
}

// O(1)
func (repository UserRepository) DeleteByEmail(email string) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	login, ok := repository.emails[email]
	if !ok {
		return fmt.Errorf("user with email '%s' not found", email)
	}

	delete(repository.users, login)
	delete(repository.emails, email)

	return nil
}

// modifies non-unique fields and (so far) does not change passwords, i.e. only username
func (repository UserRepository) Update(user *models.User) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	found_user, ok := repository.users[user.Login]
	if !ok {
		return fmt.Errorf("user with login '%s' not found", user.Login)
	}

	// if err := user.Validate(false); err != nil {
	// 	return err
	// }

	if found_user.Email != user.Email {
		if _, ok := repository.emails[user.Email]; !ok {
			return fmt.Errorf("user with email '%s' already exists", user.Email)
		}
		delete(repository.emails, found_user.Email)
		repository.emails[user.Email] = user.Login
	}

	found_user.Username = user.Username

	return nil
}
