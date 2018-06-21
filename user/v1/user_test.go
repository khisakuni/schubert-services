package route

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khisakuni/schubert-services/user/service"
	"github.com/urfave/negroni"
)

func TestCreateUser(t *testing.T) {
	app := negroni.Classic()
	s := service.Service{
		Engine: app,
		Router: New(),
	}
	req, err := http.NewRequest("POST", "/v1/u", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	s.Configure()
	s.Router.ServeHTTP(res, req)

	if res.Body.String() != "Hello world!" {
		t.Error("Expected 'Hello world!' but got", res.Body.String())
	}
}
