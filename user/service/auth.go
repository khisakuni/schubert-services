package service

import "golang.org/x/crypto/bcrypt"

type Auth interface {
	HashPassword(string) (string, error)
}

type Bcrypt struct{}

func (b Bcrypt) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
