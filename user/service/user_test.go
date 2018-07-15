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

	"github.com/lib/pq"
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

func sendRequest(t *testing.T, service Service, p params) (int, string) {
	req, err := newRequest("POST", "/api/v1/u", p)
	if err != nil {
		t.Error(err)
		return 0, ""
	}
	res := httptest.NewRecorder()
	service.Router.ServeHTTP(res, req)
	return res.Code, res.Body.String()
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
		code, message := sendRequest(t, s, test.params)
		if code != test.code {
			t.Errorf("Expected code %d, got %d\n", test.code, code)
		}
		if message != test.body {
			t.Errorf("Expected code %s, got %s\n", test.body, message)
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

	code, message := sendRequest(t, s, newParams("kohei@example.com", "password", "kohei"))
	if code != 201 {
		t.Errorf("Expected code 201, got %d\n", code)
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

	if message != string(expectedJson) {
		t.Errorf("Expected %v but got %v\n", string(expectedJson), message)
	}

	// Duplicate email
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("kohei@example.com", "kohei", "password").
		WillReturnError(&pq.Error{
			Code:    pq.ErrorCode("23505"),
			Message: "Duplicate email",
		})
	mock.ExpectRollback()

	code, _ = sendRequest(t, s, newParams("kohei@example.com", "password", "kohei"))
	if code != 400 {
		t.Errorf("Expected 400, got %d\n", code)
	}

	// Duplicate username
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("kohei1@example.com", "kohei", "password").
		WillReturnError(&pq.Error{
			Code:    pq.ErrorCode("23505"),
			Message: "Duplicate username",
		})
	mock.ExpectRollback()

	code, _ = sendRequest(t, s, newParams("kohei1@example.com", "password", "kohei"))
	if code != 400 {
		t.Errorf("Expected 400, got %d\n", code)
	}
}
