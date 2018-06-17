package main

import (
	"log"
	"net/http"

	"github.com/urfave/negroni"

	"github.com/khisakuni/schubert-services/user/v1"
)

func main() {
	app := negroni.Classic()
	app.UseHandler(route.New())
	log.Fatal(http.ListenAndServe(":8080", app))
}
