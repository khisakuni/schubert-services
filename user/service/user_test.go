package service

import (
	"encoding/json"
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

// func sendRequest(t *testing.T, service Service, p params) (int, string) {
// 	req, err := newRequest("POST", "/api/v1/u", p)
// 	if err != nil {
// 		t.Error(err)
// 		return 0, ""
// 	}
// 	res := httptest.NewRecorder()
// 	service.Router.ServeHTTP(res, req)
// 	return res.Code, res.Body.String()
// }

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
		DB:            db,
		Engine:        app,
		Router:        mux.NewRouter(),
		Authenticator: MockAuth{},
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
		if err != nil {
			t.Error(err)
		}
		res := sendRequest(t, s, req)
		if res.Code != test.code {
			t.Errorf("Expected code %d, got %d\n", test.code, res.Code)
		}
		message := res.Body.String()
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

	req, err := newRequest("POST", "/api/v1/u", newParams("kohei@example.com", "password", "kohei"))
	if err != nil {
		t.Error(err)
	}
	res := sendRequest(t, s, req)
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

	message := res.Body.String()
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

	req, err = newRequest("POST", "/api/v1/u", newParams("kohei@example.com", "password", "kohei"))
	res = sendRequest(t, s, req)
	if res.Code != 400 {
		t.Errorf("Expected 400, got %d\n", res.Code)
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

	req, err = newRequest("POST", "/api/v1/u", newParams("kohei1@example.com", "password", "kohei"))
	res = sendRequest(t, s, req)
	if res.Code != 400 {
		t.Errorf("Expected 400, got %d\n", res.Code)
	}
}
