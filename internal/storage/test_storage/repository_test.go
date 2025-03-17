package test_storage_test

import (
	"testing"
	"vox-server/internal/models"
	"vox-server/internal/storage/test_storage"

	"github.com/stretchr/testify/assert"
)

// in this package 'valid' = exists in the system

func TestUserRepository_Create(t *testing.T) {
	storage := test_storage.NewInMemoryStorage()

	// default case: new valid user
	user := &models.User{
		Login:    "user",
		Username: "username",
		Email:    "example@tmail.com",
		Password: "gooDPsswrA12",
	}
	assert.NoError(t, storage.Users().Create(user))
	assert.NotNil(t, user)

	// case : user whose login is already taken
	duplicateLoginUser := &models.User{
		Login:    "user",
		Username: "username",
		Email:    "eXampLEe@tmail.com",
		Password: "gooDPsswrA12",
	}
	assert.Error(t, storage.Users().Create(duplicateLoginUser))

	// case : user whose email is already taken
	duplicateEmailUser := &models.User{
		Login:    "new_user",
		Username: "username",
		Email:    "example@tmail.com",
		Password: "gooDPsswrA12",
	}
	assert.Error(t, storage.Users().Create(duplicateEmailUser))
}

func TestUserRepository_FindByLogin(t *testing.T) {
	storage := test_storage.NewInMemoryStorage()
	user := &models.User{
		Login:    "user",
		Username: "username",
		Email:    "example@tmail.com",
		Password: "gooDPsswrA12",
	}
	storage.Users().Create(user)

	// default case : find by login that exist
	found_user, err := storage.Users().FindByLogin(user.Login)
	assert.NoError(t, err)
	assert.NotNil(t, user, found_user)

	// case : find by login that does not exist
	_, err = storage.Users().FindByLogin("nonexistent")
	assert.Error(t, err)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	storage := test_storage.NewInMemoryStorage()
	user := &models.User{
		Login:    "user",
		Username: "username",
		Email:    "example@tmail.com",
		Password: "gooDPsswrA12",
	}
	storage.Users().Create(user)

	// default case : find by email that exist
	found_user, err := storage.Users().FindByEmail(user.Email)
	assert.NoError(t, err)
	assert.NotNil(t, user, found_user)

	// case : find by email that does not exist
	_, err = storage.Users().FindByEmail("nonexistent@tmail.com")
	assert.Error(t, err)
}

func TestUserRepository_DeleteByLogin(t *testing.T) {
	storage := test_storage.NewInMemoryStorage()
	user1 := &models.User{
		Login:    "abra",
		Username: "mutex",
		Email:    "mail@gmail.com",
		Password: "gooDPsswrA12",
	}
	user2 := &models.User{
		Login:    "kadabra",
		Username: "qwerty",
		Email:    "QwErTy@yandex.ru",
		Password: "abcdefg12134",
	}
	storage.Users().Create(user1)
	storage.Users().Create(user2)

	// default case : delete by valid login
	assert.NoError(t, storage.Users().DeleteByLogin(user1.Login))
	assert.NotNil(t, user1, user2)
	assert.EqualValues(t, storage.Users().Count(), 1)

	assert.NoError(t, storage.Users().DeleteByLogin(user2.Login))
	assert.NotNil(t, user1, user2)
	assert.True(t, storage.Users().IsEmpty())

	// case : delete by deleted login
	assert.Error(t, storage.Users().DeleteByLogin(user1.Login))
	assert.NotNil(t, user1, user2)
	assert.Error(t, storage.Users().DeleteByLogin(user2.Login))
	assert.NotNil(t, user1, user2)
	assert.True(t, storage.Users().IsEmpty())

	// case : delete by login that never existed
	assert.Error(t, storage.Users().DeleteByLogin("abrakadabra"))
	assert.NotNil(t, user1, user2)
	assert.True(t, storage.Users().IsEmpty())
}

func TestUserRepository_DeleteByEmail(t *testing.T) {
	storage := test_storage.NewInMemoryStorage()
	user1 := &models.User{
		Login:    "abra",
		Username: "mutex",
		Email:    "mail@gmail.com",
		Password: "gooDPsswrA12",
	}
	user2 := &models.User{
		Login:    "kadabra",
		Username: "qwerty",
		Email:    "QwErTy@yandex.ru",
		Password: "abcdefg12134",
	}
	storage.Users().Create(user1)
	storage.Users().Create(user2)

	// default case : delete by valid email
	assert.NoError(t, storage.Users().DeleteByEmail(user1.Email))
	assert.NotNil(t, user1, user2)
	assert.EqualValues(t, storage.Users().Count(), 1)

	assert.NoError(t, storage.Users().DeleteByEmail(user2.Email))
	assert.NotNil(t, user1, user2)
	assert.True(t, storage.Users().IsEmpty())

	// case : delete by deleted email
	assert.Error(t, storage.Users().DeleteByEmail(user1.Email))
	assert.NotNil(t, user1, user2)
	assert.Error(t, storage.Users().DeleteByEmail(user2.Email))
	assert.NotNil(t, user1, user2)
	assert.True(t, storage.Users().IsEmpty())

	// case : delete by email that never existed
	assert.Error(t, storage.Users().DeleteByEmail("abrakadabra@email.ro"))
	assert.NotNil(t, user1, user2)
	assert.True(t, storage.Users().IsEmpty())
}

func TestUserRepository_Update(t *testing.T) {
	storage := test_storage.NewInMemoryStorage()
	user := &models.User{
		Login:             "user",
		Username:          "username",
		Email:             "example@tmail.com",
		Password:          "gooDPsswrA12",
		EncryptedPassword: "encrypted_password",
	}
	storage.Users().Create(user)

	// default case : update user with valid data
	updated_user := &models.User{
		Login:             "user",
		Username:          "new_username",
		Email:             "example@tmail.com",
		Password:          "new_password",
		EncryptedPassword: "",
	}
	assert.NoError(t, storage.Users().Update(updated_user))

	foundUser, err := storage.Users().FindByLogin("user")
	assert.NoError(t, err)
	assert.Equal(t, "new_username", foundUser.Username)
	user.BeforeCreate() // encrypt password
	assert.Equal(t, user.EncryptedPassword, foundUser.EncryptedPassword)

	// case : update non-existent user
	nonExistentUser := &models.User{
		Login:    "non_existent",
		Username: "new_username",
		Email:    "new@example.com",
	}

	err = storage.Users().Update(nonExistentUser)
	assert.Error(t, err)
	assert.EqualError(t, err, "user with login 'non_existent' not found")

}
