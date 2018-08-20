package service

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func TestAuth(t *testing.T) {
	app := negroni.Classic()
	s := Service{
		Engine: app,
		Router: mux.NewRouter(),
	}
	s.Configure()

	tests := []testCase{
		testCase{
			body:   "",
			code:   200,
			params: authParams{Email: "kohei@example.com", Password: "pass123"},
		},
	}

	for _, test := range tests {
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

// 	return res.Code, res.Body.String()
