package models_test

import (
	"git-server/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *models.User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *models.User {
				return &models.User{
					Login:    "user",
					Username: "username",
					Email:    "username@xmail.com",
					Password: "passWORD0101",
				}
			},
			isValid: true,
		},
		{
			name: "empty login",
			u: func() *models.User {
				return &models.User{
					Login:    "",
					Username: "username",
					Email:    "username@xmail.com",
					Password: "passWORD0101",
				}
			},
			isValid: false,
		},
		{
			name: "invalid len(login)",
			u: func() *models.User {
				return &models.User{
					Login:    "loginloginloginloginlogin",
					Username: "username",
					Email:    "username@xmail.com",
					Password: "passWORD0101",
				}
			},
			isValid: false,
		},
		{
			name: "login with space-prefix",
			u: func() *models.User {
				return &models.User{
					Login:    "   login",
					Username: "username",
					Email:    "username@xmail.com",
					Password: "passWORD0101",
				}
			},
			isValid: false,
		},
		{
			name: "login with space-suffix",
			u: func() *models.User {
				return &models.User{
					Login:    "login   ",
					Username: "username",
					Email:    "username@xmail.com",
					Password: "passWORD0101",
				}
			},
			isValid: false,
		},
		{
			name: "bad email",
			u: func() *models.User {
				return &models.User{
					Login:    "login",
					Username: "username",
					Email:    "not-valid",
					Password: "passWORD0101",
				}
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *models.User {
				return &models.User{
					Login:    "login",
					Username: "username",
					Email:    "not-valid",
					Password: "",
				}
			},
			isValid: false,
		},
		// TODO : add more testcases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate(true))
			} else {
				assert.Error(t, tc.u().Validate(true))
			}
		})
	}
}
