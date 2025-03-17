package models_test

import (
	"testing"
	"vox-server/internal/models"

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
			name: "empty password [valid]",
			u: func() *models.User {
				return &models.User{
					Login:             "user",
					Username:          "username",
					Email:             "username@xmail.com",
					Password:          "",
					EncryptedPassword: "crypted1234",
				}
			},
			isValid: true,
		},
		{
			name: "empty password [invalid]",
			u: func() *models.User {
				return &models.User{
					Login:             "user",
					Username:          "username",
					Email:             "username@xmail.com",
					Password:          "",
					EncryptedPassword: "",
				}
			},
			isValid: false,
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
			name: "empty username",
			u: func() *models.User {
				return &models.User{
					Login:    "user",
					Username: "",
					Email:    "username@xmail.com",
					Password: "passWORD0101",
				}
			},
			isValid: false,
		},
		{
			name: "empty email",
			u: func() *models.User {
				return &models.User{
					Login:    "user",
					Username: "username",
					Email:    "",
					Password: "passWORD0101",
				}
			},
			isValid: false,
		},
		{
			name: "empty user",
			u: func() *models.User {
				return &models.User{
					Login:    "",
					Username: "",
					Email:    "",
					Password: "",
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
			name: "invalid login length",
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
			name: "invalid username length",
			u: func() *models.User {
				return &models.User{
					Login:    "login",
					Username: "usernameusernameusernameusername",
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
			name: "username with space-prefix",
			u: func() *models.User {
				return &models.User{
					Login:    "login",
					Username: "    username",
					Email:    "username@xmail.com",
					Password: "passWORD0101",
				}
			},
			isValid: false,
		},
		{
			name: "username with space-suffix",
			u: func() *models.User {
				return &models.User{
					Login:    "login",
					Username: "username    ",
					Email:    "username@xmail.com",
					Password: "passWORD0101",
				}
			},
			isValid: false,
		},
		{
			name: "wrong email format",
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
			name: "wrong email format",
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate(false))
			} else {
				if tc.name == "empty password [invalid]" {
					assert.Error(t, tc.u().Validate(true))
				} else {
					assert.Error(t, tc.u().Validate(false))
				}
			}
		})
	}
}
