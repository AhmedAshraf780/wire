package wire

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

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
	wg := &sync.WaitGroup{}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		wg.Add(1)
		go app.handleConnection(conn, wg)
	}
}

func (app *Application) handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		// request line
		request := Request[[]byte]{
			Headers: make(map[string]string),
		}

		line, err := reader.ReadString('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			return
		}
		tokens, err := utils.ParseRequestLine(line)
		if err != nil {
			// TODO: send http response
			return
		}

		request.Method = tokens[0]
		request.Path = tokens[1]
		request.Version = tokens[2]
		// here we can check if we have a handler for this path
		path, query, found := strings.Cut(tokens[1], "?")
		if found {
			request.Query = ParseQuery(query)
		}
		key := utils.GenerateHandlerKey(request.Method, path)
		_, ok := app.staticRoutes[key]
		d := -1
		orgpath := request.Path
		if !ok {
			// check in dynamic routes
			params, org, idx, meth := CheckDynamicPath(app.dynamicRoutes, path)
			if idx != -1 && meth == request.Method {
				d = idx
				request.Params = params
				orgpath = org
				goto label1
			}
			resp := utils.MakeResponse(http.StatusMethodNotAllowed, "No method exists", []byte("Method or path not found"), map[string]string{}, request.Version)
			_, err := conn.Write(resp)
			if err != nil {
				log.Println(err)
			}
			return
		}
	label1:
		if d != -1 {
			_ = app.dynamicRoutes[d].Handler
		}
		// request headers
		for {
			line, err = reader.ReadString('\n')
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Println(err)
				return
			}
			if line == "\r\n" {
				break
			}
			key, value, err := utils.ParseHeader(line)
			if err != nil {
				log.Println(err)
				// TODO: send http response
				return
			}
			request.Headers[key] = value
		}
		// request body
		lengthStr, _ := request.Headers["Content-Length"]
		length, _ := strconv.Atoi(lengthStr)
		body := make([]byte, length)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			log.Println(err)
			return
		}
		request.Body = body
		request.Context = make(map[string]interface{})
		// now the request is ready
		app.handleRequest(request, orgpath, conn)
	}
}

func (app *Application) handleRequest(request Request[[]byte], path string, conn net.Conn) {
	// run global middlewares first
	for _, mid := range app.globalMiddlewares {
		err := mid.Handle(&request, &Response[any]{
			Headers: make(map[string]string),
		}, conn)
		if err != nil && errors.Is(err, ErrNext) {
			continue
		} else {
			return
		}
	}

	middlewares, ok := app.middlewares[path]
	if !ok {
		return
	}
	for _, middleware := range middlewares {
		err := middleware.Handle(&request, &Response[any]{
			Headers: make(map[string]string),
		}, conn)
		if err != nil && errors.Is(err, ErrNext) {
			continue
		} else {
			return
		}
	}
}
