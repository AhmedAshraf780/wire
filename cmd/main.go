package main

import (
	"httpserver/internals/wire"
	"net/http"
)

type User struct {
	Name  string
	Email string
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

var Users []User

func main() {
	app := wire.NewApplication()
	wire.GET(app, "/", Healthz)
	wire.POST(app, "/api/v1/users", CreateUser)
	wire.GET(app, "/api/v2/users", GetAllUsers)
	wire.GET(app, "/api/v2/users/:id/email/:email", GetUserByID)
	app.Listen(3000)
}

func Healthz(req wire.Request[wire.EmptyBody], resp *wire.Response[responseType]) {
	resp.Write(http.StatusOK, responseType{Ok: true, Message: "server is ok", Users: nil})
}

func CreateUser(req wire.Request[requestType], resp *wire.Response[responseType]) {
	name := req.Body.Name
	email := req.Body.Email

	if name == "" || email == "" {
		resp.Write(http.StatusBadRequest, responseType{Ok: false, Message: "name or email required", Users: nil})
		return
	}
	Users = append(Users, User{name, email})
	resp.Write(http.StatusOK, responseType{Ok: true, Message: "ok", Users: Users})
	return
}

func GetAllUsers(req wire.Request[wire.EmptyBody], resp *wire.Response[responseType]) {
	resp.Write(http.StatusOK, responseType{Ok: true, Message: "server is ok", Users: Users})
}

func GetUserByID(req wire.Request[wire.EmptyBody], resp *wire.Response[responseType]) {
	id, ok := req.Params["id"]
	name, _ := req.Query["name"]
	if !ok {
		resp.Write(http.StatusBadRequest, responseType{Ok: false, Message: "id required", Users: nil})
		return
	}
	for _, user := range Users {
		if user.Name == name {
			resp.Write(http.StatusOK, responseType{Ok: true, Message: "ok", Users: Users})
			return
		}
		if user.Email == id {
			resp.Write(http.StatusOK, responseType{Ok: false, Message: "ok dude", Users: Users})
			return
		}
	}
	resp.Write(http.StatusNotFound, responseType{Ok: false, Message: "user not found", Users: nil})
}
