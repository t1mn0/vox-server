package sqlstorage_test

import (
	"git-server/internal/models"
	"git-server/internal/storage/sqlstorage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstorage.TestDB(t, databaseURL)
	st := sqlstorage.New(db)
	defer teardown("users")

	user := &models.User{
		Login:    "user",
		Username: "username",
		Email:    "example@tmail.com",
		Password: "gooDPsswrA12",
	}

	err := st.User().Create(user)

	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstorage.TestDB(t, databaseURL)
	st := sqlstorage.New(db)
	defer teardown("users")

	login := "user"
	_, err := st.User().FindByLogin(login)

	assert.Error(t, err)

	st.User().Create(&models.User{
		Login:    "user",
		Username: "username",
		Email:    "example@tmail.com",
		Password: "gooDPsswrA12",
	})

	u, err := st.User().FindByLogin(login)

	assert.NoError(t, err)
	assert.NotNil(t, u)

}
