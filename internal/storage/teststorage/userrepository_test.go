package teststorage_test

import (
	"git-server/internal/models"
	"git-server/internal/storage/teststorage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	st := teststorage.New()
	user := &models.User{
		Login:    "user",
		Username: "username",
		Email:    "example@tmail.com",
		Password: "gooDPsswrA12",
	}
	assert.NoError(t, st.User().Create(user))
	assert.NotNil(t, user)
}

func TestUserRepository_FindByLogin(t *testing.T) {
	st := teststorage.New()
	user := &models.User{
		Login:    "user",
		Username: "username",
		Email:    "example@tmail.com",
		Password: "gooDPsswrA12",
	}
	st.User().Create(user)
	_, err := st.User().FindByLogin(user.Login)

	assert.NoError(t, err)
	assert.NotNil(t, user)
}
