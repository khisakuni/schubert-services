package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	Username        string `json:"username"`
}

func (s *Service) createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.Write([]byte("oh no!"))
		return
	}

	//hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	hashedPass, err := s.Auth.HashPassword(user.Password)
	if err != nil {
		w.Write([]byte("oh no!"))
		return
	}

	tx, err := s.DB.Begin()
	if err != nil {
		w.Write([]byte("oh no!"))
		return
	}

	sql := `
		INSERT INTO users (email, username, password) VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(sql, user.Email, user.Username, hashedPass)
	if err != nil {
		fmt.Printf("err > %v\n", err)
		w.Write([]byte("oh no!"))
		return
	}

	switch err {
	case nil:
		err = tx.Commit()
	default:
		tx.Rollback()
	}

	w.Write([]byte(user.Email))
}
