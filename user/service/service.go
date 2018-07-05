package service

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

type Service struct {
	DB     *sql.DB
	Engine *negroni.Negroni
	Port   string
	Router *mux.Router
	Auth   Auth
}

func (s *Service) Configure() {
	//s.Engine.Use(negroni.HandlerFunc(middleware.WithDB(s.DB)))
	s.Router.HandleFunc("/api/v1/u", s.createUser).Methods("POST")
	s.Engine.UseHandler(s.Router)
}

func (s *Service) Run() {
	log.Fatal(http.ListenAndServe(s.Port, s.Engine))
}
