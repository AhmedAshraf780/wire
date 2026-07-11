package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/AhmedAshraf780/wire/internals/wire"
)

type User struct {
	Name  string
	Email string
}

type responseType struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Users   []User `json:"users"`
}

var Users []User

func main() {
	// wire
	app := wire.NewApplication()
	wire.POST(app, "/users", createUser)
	app.Listen(3000)
	//http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
	//	var u User
	//	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
	//		http.Error(w, err.Error(), http.StatusBadRequest)
	//	}
	//	Users = append(Users, u)
	//	w.Write([]byte("Users were added successfully"))
	//})
	//http.ListenAndServe(":3000", nil)
}

func createUser(req *wire.Request[User], resp *wire.Response[responseType]) error {
	Users = append(Users, User{
		Name:  req.Body.Name,
		Email: req.Body.Email,
	})
	return resp.Write(http.StatusCreated, responseType{
		Ok:      true,
		Message: "user created",
		Users:   Users,
	})
}
