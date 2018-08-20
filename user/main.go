package main

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/khisakuni/schubert-services/user/database"
	"github.com/khisakuni/schubert-services/user/service"
)

func main() {
	app := negroni.Classic()
	db, err := database.New()
	if err != nil {
		panic(err)
	}

	service := service.Service{
		Authenticator: service.Bcrypt{},
		DB:            db,
		Engine:        app,
		Port:          ":8080",
		Router:        mux.NewRouter(),
	}

	service.Configure()
	fmt.Println("Listening on port 8080")
	service.Run()
}
