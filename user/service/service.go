package service

import (
	"database/sql"

	"fmt"
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

type handlerError struct {
	code    int
	message string
}

func (he *handlerError) Error() string {
	return fmt.Sprintf("%d - %s", he.code, he.message)
}

type handler func(w http.ResponseWriter, r *http.Request) error

func handleError(handlerFunc handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handlerFunc(w, r)
		if err != nil {
			if he, ok := err.(*handlerError); ok {
				w.WriteHeader(he.code)
				w.Write([]byte(he.message))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Something went wrong - %s", err)))
		}
	}
}

func (s *Service) Configure() {
	s.Router.HandleFunc("/api/v1/u", handleError(s.createUser)).Methods("POST")
	s.Engine.UseHandler(s.Router)
}

func (s *Service) Run() {
	log.Fatal(http.ListenAndServe(s.Port, s.Engine))
}
