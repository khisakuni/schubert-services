package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID              int64  `json:"id,omitempty"`
	Email           string `json:"email"`
	Password        string `json:"password,omitempty"`
	PasswordConfirm string `json:"passwordConfirm,omitempty"`
	Username        string `json:"username"`
}

const passwordMinLen = 7

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

	hashedPass, err := s.Auth.HashPassword(user.Password)
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	sql := `
		INSERT INTO users (email, username, password) VALUES ($1, $2, $3)
	`
	result, err := tx.Exec(sql, user.Email, user.Username, hashedPass)
	if err != nil {
		return err
	}

	switch err {
	case nil:
		err = tx.Commit()
	default:
		tx.Rollback()
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	jsonRes, err := json.Marshal(User{
		ID:       ID,
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
