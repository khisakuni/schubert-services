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

// Need to test:
// - invalid requests
//   -> response
//   -> no DB
// - valid requests
//   -> response
//   -> DB

type params struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	Username        string `json:"username"`
}

func newParams(email, password, username string, passwordConfirm ...string) params {
	var confirm string
	if len(passwordConfirm) <= 0 {
		confirm = password
	} else {
		confirm = passwordConfirm[0]
	}
	return params{
		Email:           email,
		Password:        password,
		Username:        username,
		PasswordConfirm: confirm,
	}
}

func newRequest(method, endpoint string, p params) (*http.Request, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(p)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(method, endpoint, body)
}

func TestCreateUser(t *testing.T) {
	app := negroni.Classic()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()

	// Initialize service
	s := Service{
		DB:     db,
		Engine: app,
		Router: mux.NewRouter(),
		Auth:   MockAuth{},
	}
	s.Configure()

	// Test invalid requests
	// Missing email
	p := newParams("", "password", "kohei")
	req, err := newRequest("POST", "/api/v1/u", p)
	res := httptest.NewRecorder()
	s.Router.ServeHTTP(res, req)
	if res.Code != 400 {
		t.Errorf("Expected code 400, got %d\n", res.Code)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	// Missing password
	p = newParams("kohei@example.com", "", "kohei")
	req, err = newRequest("POST", "/api/v1/u", p)
	res = httptest.NewRecorder()
	s.Router.ServeHTTP(res, req)
	if res.Code != 400 {
		t.Errorf("Expected code 400, got %d\n", res.Code)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	// Missing username
	p = newParams("kohei@example.com", "password", "")
	req, err = newRequest("POST", "/api/v1/u", p)
	res = httptest.NewRecorder()
	s.Router.ServeHTTP(res, req)
	if res.Code != 400 {
		t.Errorf("Expected code 400, got %d\n", res.Code)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	// Mismatches passwordConfirm
	p = newParams("kohei@example.com", "password", "kohei", "wrong")
	req, err = newRequest("POST", "/api/v1/u", p)
	res = httptest.NewRecorder()
	s.Router.ServeHTTP(res, req)
	if res.Code != 400 {
		t.Errorf("Expected code 400, got %d\n", res.Code)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	// Correct params
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users").WithArgs("kohei@example.com", "kohei", "password").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	p = newParams("kohei@example.com", "password", "kohei")
	req, err = newRequest("POST", "/api/v1/u", p)
	res = httptest.NewRecorder()
	s.Router.ServeHTTP(res, req)
	if res.Code != 201 {
		t.Errorf("Expected code 201, got %d\n", res.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	expectedUser := User{
		ID:       1,
		Email:    "kohei@example.com",
		Username: "kohei",
	}

	expectedJson, err := json.Marshal(expectedUser)
	if err != nil {
		t.Error(err)
	}

	if res.Body.String() != string(expectedJson) {
		t.Errorf("Expected %v but got %v\n", string(expectedJson), res.Body.String())
	}

}
