package main

import (
	"fmt"
	"httpserver/internals/wire"
)

func main() {
	app := wire.NewApplication()
	app.GET("/hello", HandleHello)
	app.Listen(3000)
}

func HandleHello(req wire.Request, resp *wire.Response) {
	resp.Body = []byte("Hello World from handler")
	resp.StatusCode = 200
	resp.Headers["Content-Type"] = "text/plain; charset=utf-8"
	resp.Headers["Content-Length"] = fmt.Sprintf("%d", len(resp.Body))
	resp.Headers["host"] = "localhost"
}
