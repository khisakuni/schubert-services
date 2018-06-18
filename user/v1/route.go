package route

import (
	"github.com/gorilla/mux"
)

// New returns Router.
func New() *mux.Router {
	router := mux.NewRouter()

	// Add routes
	router.HandleFunc("/v1/u", createUser).Methods("POST")

	return router
}
