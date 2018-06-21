package main

import (
	"github.com/urfave/negroni"

	"github.com/khisakuni/schubert-services/user/database"
	"github.com/khisakuni/schubert-services/user/service"
	"github.com/khisakuni/schubert-services/user/v1"
)

func main() {
	app := negroni.Classic()
	db, err := database.New()
	if err != nil {
		panic(err)
	}

	service := service.Service{
		DB:     db,
		Engine: app,
		Port:   ":8080",
		Router: route.New(),
	}
	service.Run()
}
