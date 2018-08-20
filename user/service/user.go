package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lib/pq"
)

type User struct {
	ID              int    `json:"id,omitempty"`
	Email           string `json:"email"`
	Password        string `json:"password,omitempty"`
	PasswordConfirm string `json:"passwordConfirm,omitempty"`
	Username        string `json:"username"`
}

const (
	passwordMinLen = 7
	duplicateError = pq.ErrorCode("23505")
)

// TODO:
//   - Add more rubst email validation

func (u *User) validate() error {
	var message string
	if len(u.Email) <= 0 {
		message = "Missing email"
	}
	if len(u.Username) <= 0 {
		message = "Missing username"
	}
	if len(u.Password) <= passwordMinLen {
		message = fmt.Sprintf("Password must be at least %d characters", passwordMinLen+1)
	}
	if u.Password != u.PasswordConfirm {
		message = "Passwords don't match"
	}
	if message != "" {
		return &handlerError{
			code:    http.StatusBadRequest,
			message: message,
		}
	}
	return nil
}

func (s *Service) createUser(w http.ResponseWriter, r *http.Request) error {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return err
	}

	if err := user.validate(); err != nil {
		return err
	}

	hashedPass, err := s.Authenticator.HashPassword(user.Password)
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	var id int
	sql := `
		INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING id
	`
	err = tx.QueryRow(sql, user.Email, user.Username, hashedPass).Scan(&id)
	switch err {
	case nil:
		err = tx.Commit()
	default:
		tx.Rollback()
	}

	if err, ok := err.(*pq.Error); ok {
		if err.Code == duplicateError {
			return &handlerError{
				code:    http.StatusBadRequest,
				message: err.Message,
			}
		}
	}

	jsonRes, err := json.Marshal(User{
		ID:       id,
		Email:    user.Email,
		Username: user.Username,
	})
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonRes)
	return nil
}
