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
}

// 	return res.Code, res.Body.String()
