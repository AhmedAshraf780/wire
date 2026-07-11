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

func (app *Application) Listen(port int) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	log.Println("Listening on " + addr)
	jobs := make(chan net.Conn, 1000)

	for i := 0; i < 200; i++ {
		go func() {
			for conn := range jobs {
				app.handleConnection(conn)
			}
		}()
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			println("CAN'T ACCEPT:", err)
			continue
		}
		jobs <- conn // blocks when pool is saturated
	}
}

func (app *Application) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		conn.SetReadDeadline(time.Now().Add(40 * time.Second))
		conn.SetWriteDeadline(time.Now().Add(40 * time.Second))
		// now the request is ready
		request, orgPath, ok := ReadAndParseRequest(app, reader, conn)
		if !ok {
			return
		}
		if request == nil {
			continue
		}
		app.handleRequest(request, orgPath, conn)
	}
}

func (app *Application) handleRequest(request *Request[[]byte], path string, conn net.Conn) {
	// run global middlewares first
	resp := &Response[any]{
		Headers: make(map[string]string),
	}
	for _, mid := range app.globalMiddlewares {
		err := mid.Handle(request, resp, conn)
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
		err := middleware.Handle(request, resp, conn)
		if err != nil && errors.Is(err, ErrNext) {
			continue
		}
		if err != nil {
			return
		}
	}
}
