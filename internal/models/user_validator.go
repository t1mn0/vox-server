package models

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

var validate *validator.Validate

// if password_check == true => field 'password_check' must be filled
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
