package auth

import (
	"golang.org/x/crypto/bcrypt"
	"errors"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost + 2)
	if err != nil {
		return "", errors.New("Error hashing dat password")
	}
	return string(hash), nil
}

func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
