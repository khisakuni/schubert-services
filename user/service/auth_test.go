package service

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestAuth(t *testing.T) {
	app := negroni.Classic()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	authenticator := Bcrypt{}
	s := Service{
		Authenticator: authenticator,
		Engine:        app,
		Router:        mux.NewRouter(),
		DB:            db,
	}
	s.Configure()

	email := "kohei@example.com"
	password := "pass123"
	hashed, err := s.Authenticator.HashPassword(password)
	if err != nil {
		t.Error(err)
		return
	}

	tests := []testCase{
		testCase{
			before: func(t *testing.T) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT password FROM users WHERE").
					WithArgs(email).
					WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow(hashed))
			},
			body:   "",
			code:   200,
			params: authParams{Email: email, Password: password},
		},
		testCase{
			before: func(t *testing.T) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT password FROM users WHERE").
					WithArgs(email).
					WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow(hashed))
			},
			body:   "Unauthorized",
			code:   401,
			params: authParams{Email: email, Password: password + "wrong"},
		},
		testCase{
			before: func(t *testing.T) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT password FROM users WHERE").
					WithArgs("kohei@example.net").
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			body:   "Not found",
			code:   404,
			params: authParams{Email: "kohei@example.net", Password: password},
		},
	}

	for _, test := range tests {
		if test.before != nil {
			test.before(t)
		}
		req, err := newRequest("POST", "/api/v1/auth", test.params)
		if err != nil {
			t.Error(err)
		}
		res := sendRequest(t, s, req)
		if res.Code != test.code {
			t.Errorf("Expected %d, got %d", test.code, res.Code)
		}

		body := res.Body.String()
		if body != test.body {
			t.Errorf("Expected %s, got %s", test.body, body)
		}
	}
}
