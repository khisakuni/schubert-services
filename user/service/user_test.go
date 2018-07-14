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

type testCase struct {
	params
	code int
	body string
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

	tests := []testCase{
		testCase{
			params: newParams("", "password", "kohei"),
			code:   400,
			body:   "Missing email",
		},
		testCase{
			params: newParams("kohei@example.com", "", "kohei"),
			code:   400,
			body:   "Password must be at least 8 characters",
		},
		testCase{
			params: newParams("kohei@example.com", "password", ""),
			code:   400,
			body:   "Missing username",
		},
		testCase{
			params: newParams("kohei@example.com", "password", "kohei", "wrong"),
			code:   400,
			body:   "Passwords don't match",
		},
		testCase{
			params: newParams("kohei@example.com", "pass", "kohei"),
			code:   400,
			body:   "Password must be at least 8 characters",
		},
	}

	for _, test := range tests {
		req, err := newRequest("POST", "/api/v1/u", test.params)
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		if res.Code != test.code {
			t.Errorf("Expected code %d, got %d\n", test.code, res.Code)
		}
		if res.Body.String() != test.body {
			t.Errorf("Expected code %s, got %s\n", test.body, res.Body.String())
		}
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	}

	// Correct params
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("kohei@example.com", "kohei", "password").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	p := newParams("kohei@example.com", "password", "kohei")
	req, err := newRequest("POST", "/api/v1/u", p)
	res := httptest.NewRecorder()
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
