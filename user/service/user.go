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

func (s *Service) createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.Write([]byte("oh no!"))
		return
	}

	// TODO: Robustify
	if len(user.Email) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing email"))
		return
	}

	// TODO: Are there other password requirements?
	if len(user.Password) <= 7 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Password must be at least 8 characters"))
		return
	}

	if user.Password != user.PasswordConfirm {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Passwords don't match"))
		return
	}

	if len(user.Username) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing username"))
		return
	}

	//hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	hashedPass, err := s.Auth.HashPassword(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oh no!"))
		return
	}

	tx, err := s.DB.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oh no!"))
		return
	}

	sql := `
		INSERT INTO users (email, username, password) VALUES ($1, $2, $3)
	`
	result, err := tx.Exec(sql, user.Email, user.Username, hashedPass)
	if err != nil {
		w.Write([]byte("oh no!"))
		return
	}

	switch err {
	case nil:
		err = tx.Commit()
	default:
		tx.Rollback()
	}

	ID, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oh no!"))
		return
	}

	jsonRes, err := json.Marshal(User{
		ID:       ID,
		Email:    user.Email,
		Username: user.Username,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oh no!"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonRes)
}
