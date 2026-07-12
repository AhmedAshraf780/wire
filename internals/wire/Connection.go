package wire

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AhmedAshraf780/wire/internals/utils"
)

type Client struct {
	conn   net.Conn
	closed bool
}

func (app *Application) Listen(port int) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		println("Couldn't accept")
		return
	}

	defer listener.Close()

	log.Println("Listening on " + addr)

	i := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			println("CAN'T ACCEPT:", err)
			continue
		}

		go app.handleConnection(Client{conn, false})
		println("Accepting connection ", i)
		i++
	}
}

func (app *Application) handleConnection(client Client) {
	defer func() {
		client.conn.Close()
	}()
	reader := bufio.NewReader(client.conn)
	for {
		client.conn.SetDeadline(time.Now().Add(30 * time.Second))
		request, orgPath, ok := readAndParseRequest(app, reader, client)
		if !ok {
			client.closed = true
			return
		}
		if request == nil {
			continue
		}
		app.handleRequest(request, orgPath, client)
		if client.closed {
			return
		}
	}
}

func (app *Application) handleRequest(request *Request[[]byte], path string, client Client) {
	// run global middlewares first
	resp := &Response[any]{
		Headers: make(map[string][]string),
	}
	for _, mid := range app.globalMiddlewares {
		err := mid.Handle(request, resp, client)
		if err != nil && errors.Is(err, ErrNext) {
			continue
		}
		if err != nil {
			return
		}
	}

	key := utils.GenerateHandlerKey(request.Method, path)
	middlewares, ok := app.middlewares[key]
	if !ok {
		return
	}
	for _, middleware := range middlewares {
		err := middleware.Handle(request, resp, client)
		if err != nil && errors.Is(err, ErrNext) {
			continue
		}
		if err != nil {
			return
		}
	}
}
