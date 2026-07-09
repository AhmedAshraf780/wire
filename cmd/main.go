package main

import (
	"httpserver/internals/wire"
)

type User struct {
	Name  string
	Email string
}

var Users []User

func main() {
	app := wire.NewApplication()
	wire.GET(app, "/", Healthz)
	wire.POST(app, "/api/v1/users", CreateUser)
	wire.GET(app, "/api/v1/users", GetAllUsers)
	app.Listen(3000)
}

type requestType struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type responseType struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Users   []User `json:"users"`
}

func Healthz(req wire.Request[any], resp *wire.Response[responseType]) {
	resp.Write(200, responseType{Ok: true, Message: "server is ok", Users: nil})
}

func CreateUser(req wire.Request[requestType], resp *wire.Response[responseType]) {
	name := req.Body.Name
	email := req.Body.Email
	if name == "" || email == "" {
		resp.Write(400, responseType{false, "Name or email is invalid", nil})
		return
	}

	Users = append(Users, User{name, email})
	resp.Write(201, responseType{true, "User created successfully", nil})
	return
}

func GetAllUsers(req wire.Request[any], resp *wire.Response[responseType]) {
	resp.Write(200, responseType{true, "All users successfully", Users})
}
