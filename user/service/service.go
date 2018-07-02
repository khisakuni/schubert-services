package service

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"log"
	"net/http"

	"github.com/khisakuni/schubert-services/user/middleware"
)

type Service struct {
	DB     *sql.DB
	Engine *negroni.Negroni
	Port   string
	Router *mux.Router
}

func (s *Service) Configure() {
	s.Engine.Use(negroni.HandlerFunc(middleware.WithDB(s.DB)))
}

func (s *Service) Run() {
	s.Configure()
	s.Engine.UseHandler(s.Router)
	log.Fatal(http.ListenAndServe(s.Port, s.Engine))
}