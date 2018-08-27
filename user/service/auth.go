package service

import (
	"encoding/json"
	"net/http"

	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type Authenticator interface {
	HashPassword(string) (string, error)
	Compare(hashed, password string) error
}

type authParams struct {
	Email    string
	Password string
}

func (s *Service) auth(w http.ResponseWriter, r *http.Request) error {
	var p authParams
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	var password string
	stmt := `SELECT password FROM users WHERE email = ?`
	err = tx.QueryRow(stmt, p.Email).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			return &handlerError{
				code:    http.StatusNotFound,
				message: "Not found",
			}
		}
		return err
	}

	err = s.Authenticator.Compare(password, p.Password)
	if err != nil {
		return &handlerError{
			code:    http.StatusUnauthorized,
			message: "Unauthorized",
		}
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

func (b Bcrypt) Compare(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
