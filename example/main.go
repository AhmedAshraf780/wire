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
	app := wire.NewApplication()
	wire.GET(app, "/users", getUsers)
	app.Listen(4000)
}
func getUsers(req *wire.Request[wire.EmptyBody], resp *wire.Response[responseType]) error {
	return resp.Write(http.StatusOK, responseType{Ok: true, Users: []User{{Name: "ahmed", Email: "ashraf"}}})
}
