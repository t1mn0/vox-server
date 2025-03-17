package models

import "golang.org/x/crypto/bcrypt"

func encryptString(str string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

type User struct {
	Login             string `validate:"required,lte=20" json:"login"`
	Username          string `validate:"required,lte=20" json:"username"`
	Email             string `validate:"required,email" json:"email"`
	Password          string `validate:"required,min=8,max=40" json:"password,omitempty"`
	EncryptedPassword string `validate:"omitempty" json:"-"`
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

func (u *User) Sanitize() {
	u.Password = ""
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}
