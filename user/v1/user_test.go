package route

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khisakuni/schubert-services/user/service"
	"github.com/urfave/negroni"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreateUser(t *testing.T) {
	app := negroni.Classic()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()

	s := service.Service{
		DB:     db,
		Engine: app,
		Router: New(),
	}
	req, err := http.NewRequest("POST", "/v1/u", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	res := httptest.NewRecorder()
	s.Configure()
	s.Router.ServeHTTP(res, req)

	if res.Body.String() != "Hello world!" {
		t.Error("Expected 'Hello world!' but got", res.Body.String())
	}
}
