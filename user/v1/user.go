package route

import (
	"encoding/json"
	"net/http"
)

type User struct {
	email           string
	password        string
	passwordConfirm string
	username        string
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.Write([]byte("oh no!"))
		return
	}
	w.Write([]byte(user.email))
}
