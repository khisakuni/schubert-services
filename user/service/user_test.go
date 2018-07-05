package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type MockAuth struct{}

func (a MockAuth) HashPassword(password string) (string, error) {
	return password, nil
}

func TestCreateUser(t *testing.T) {
	app := negroni.Classic()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()

	s := Service{
		DB:     db,
		Engine: app,
		Router: mux.NewRouter(),
		Auth:   MockAuth{},
	}

	type params struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"passwordConfirm"`
		Username        string `json:"username"`
	}

	email := "kohei@example.com"
	pass := "password123"
	username := "kohei"
	body := new(bytes.Buffer)
	p := params{
		Email:           email,
		Password:        pass,
		PasswordConfirm: pass,
		Username:        username,
	}
	err = json.NewEncoder(body).Encode(p)
	if err != nil {
		t.Fatal(err)
		return
	}

	req, err := http.NewRequest("POST", "/api/v1/u", body)
	if err != nil {
		t.Fatal(err)
		return
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users").WithArgs(email, username, pass).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res := httptest.NewRecorder()
	s.Configure()
	s.Router.ServeHTTP(res, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if res.Body.String() != "Hello world!" {
		t.Error("Expected 'Hello world!' but got", res.Body.String())
	}

}
