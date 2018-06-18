package main

import (
	"log"
	"net/http"

	"github.com/urfave/negroni"

	"github.com/khisakuni/schubert-services/user/database"
	"github.com/khisakuni/schubert-services/user/middleware"
	"github.com/khisakuni/schubert-services/user/v1"
)

func main() {
	app := negroni.Classic()
	db, err := database.New()
	if err != nil {
		panic(err)
	}
	app.Use(negroni.HandlerFunc(middleware.WithDB(db)))

	app.UseHandler(route.New())
	log.Fatal(http.ListenAndServe(":8080", app))
}
