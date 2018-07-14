package service

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID              int64  `json:"id,omitempty"`
	Email           string `json:"email"`
	Password        string `json:"password,omitempty"`
	PasswordConfirm string `json:"passwordConfirm,omitempty"`
	Username        string `json:"username"`
}

// TODO: Implement better error handling

func (s *Service) createUser(w http.ResponseWriter, r *http.Request) error {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return err
	}

	// TODO: Robustify
	if len(user.Email) <= 0 {
		return &handlerError{
			code:    http.StatusBadRequest,
			message: "Missing email",
		}
	}

	// TODO: Are there other password requirements?
	if len(user.Password) <= 7 {
		return &handlerError{
			code:    http.StatusBadRequest,
			message: "Password must be at least 8 characters",
		}
	}

	if user.Password != user.PasswordConfirm {
		return &handlerError{
			code:    http.StatusBadRequest,
			message: "Passwords don't match",
		}
	}

	if len(user.Username) <= 0 {
		return &handlerError{
			code:    http.StatusBadRequest,
			message: "Missing username",
		}
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
