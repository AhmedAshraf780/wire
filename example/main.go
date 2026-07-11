package main

import (
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

func main() {
	app := wire.NewApplication()
	wire.USE(app, auth)
	wire.GET(app, "/users/:id/emails/:email", first, Second)
	app.Listen(3000)
}

func first(req *wire.Request[wire.EmptyBody], resp *wire.Response[User]) error {
	req.Context["User"] = "Hend roshdy"
	return wire.ErrNext
}

func auth(req *wire.Request[wire.EmptyBody], resp *wire.Response[User]) error {
	req.Context["token"] = "verysecrettoken"
	return wire.ErrNext
}

//	func third(req *wire.Request[wire.EmptyBody], resp *wire.Response[User]) error {
//		user := req.Context["token"].(string)
//		req.Context["Email"] = user + "@gmail.com"
//		return wire.ErrNext
//	}

func Second(req *wire.Request[wire.EmptyBody], resp *wire.Response[User]) error {
	token, ok := req.Context["token"].(string)
	if !ok {
		return resp.Write(500, User{Name: "undefined", Email: "undefined"})
	}
	id, ok := req.Params["id"]
	if !ok {
		return resp.Write(500, User{Name: "You missed id param", Email: "undefined"})
	}
	email, ok := req.Params["email"]
	if !ok {
		return resp.Write(500, User{Name: "You missed email param", Email: "undefined"})
	}

	user := User{
		Name:  id + token,
		Email: email,
	}

	b, ok := req.Query["bool"]
	if !ok {
		return resp.Write(500, User{Name: "You missed bool Query", Email: "undefined"})
	}
	println(b)
	return resp.Write(200, user)
}
