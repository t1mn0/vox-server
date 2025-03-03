package models

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Login             string `validate:"required,lte=20"`
	Username          string `validate:"required,lte=20"`
	Email             string `validate:"required,email"`
	Password          string `validate:"omitempty,min=8,max=40"`
	EncryptedPassword string `validate:"omitempty"`
}

var validate *validator.Validate

func (u *User) Validate(password_check bool) error {
	if validate == nil {
		validate = validator.New()
	}

	if strings.HasPrefix(u.Login, " ") || strings.HasSuffix(u.Login, " ") {
		return fmt.Errorf("login should not start or end with spaces")
	}
	if strings.HasPrefix(u.Username, " ") || strings.HasSuffix(u.Username, " ") {
		return fmt.Errorf("username should not start or end with spaces")
	}

	if hasSpecialCharacters(&u.Login) {
		return fmt.Errorf("login contains special characters")
	}
	if hasSpecialCharacters(&u.Username) {
		return fmt.Errorf("username contains special characters")
	}

	var err error

	if password_check {
		err = validate.StructExcept(u, "EncryptedPassword")
	} else {
		err = validate.StructExcept(u, "Password")
	}

	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	return nil
}

func hasSpecialCharacters(s *string) bool {
	for _, r := range *s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return true
		}
	}
	return false
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}

	return nil
}

func encryptString(str string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
