package service

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Authenticator interface {
	HashPassword(string) (string, error)
}

type authParams struct {
	Email    string
	Password string
}

func (s *Service) auth(w http.ResponseWriter, r *http.Request) error {
	var p authParams
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

type Bcrypt struct{}

func (b Bcrypt) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
