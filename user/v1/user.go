package route

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	Username        string `json:"username"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.Write([]byte("oh no!"))
		return
	}
	w.Write([]byte(user.Email))
}
