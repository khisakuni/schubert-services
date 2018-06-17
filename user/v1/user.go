package route

import (
	"fmt"
	"net/http"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "creating user!")
}
