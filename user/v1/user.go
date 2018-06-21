package route

import (
	"net/http"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}
