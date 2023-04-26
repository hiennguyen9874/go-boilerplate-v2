package cryptpass

import (
	"errors"

	"github.com/hiennguyen9874/go-boilerplate/pkg/httpErrors"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", httpErrors.ErrBadRequest(errors.New("can not hash password"))
	}
	return string(hashedPassword), nil
}

func ComparePassword(password string, hashedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}
