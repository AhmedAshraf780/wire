package wire

import (
	"bufio"
	"fmt"
	"httpserver/internals/utils"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Route struct {
	Method   string
	Pattern  string
	Segments []string
	Handler  Handler
}
type Application struct {
	staticRoutes  map[string]Handler
	dynamicRoutes []Route
}

func NewApplication() *Application {
	return &Application{
		staticRoutes: make(map[string]Handler),
	}
}

func (app *Application) Listen(port int) {
	// printing static and dynamic routes
	for _, route := range app.dynamicRoutes {
		for _, segment := range route.Segments {
			print(segment)
		}
	}
	for p, _ := range app.staticRoutes {
		println(p)
	}
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
		key := utils.GenerateHandlerKey(request.Method, request.Path)
		handler, ok := app.staticRoutes[key]
		d := -1
		if !ok {
			// we send an http error response
			// check in dynamic routes
			params, q, idx, meth := CheckDynamicPath(app.dynamicRoutes, request.Path)
			if idx != -1 && meth == request.Method {
				d = idx
				request.Query = q
				request.Params = params
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
			handler = app.dynamicRoutes[d].Handler
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
		// here we need to check the handler mode
		err = handler.Handle(request, &Response[any]{
			Headers: map[string]string{},
		}, conn)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func GET[T, V any](app *Application, path string, callback func(Request[T], *Response[V])) {
	handler := &WireHandler[T, V]{
		Path:     path,
		Callback: callback,
	}
	if utils.StaticPath(path) {
		app.staticRoutes["GET:"+path] = handler
	} else {
		app.dynamicRoutes = append(app.dynamicRoutes, Route{
			Method:   "GET",
			Pattern:  path,
			Segments: strings.Split(path, "/"),
			Handler:  handler,
		})
	}
}

func POST[T, V any](app *Application, path string, callback func(Request[T], *Response[V])) {
	handler := &WireHandler[T, V]{
		Path:     path,
		Callback: callback,
	}
	if utils.StaticPath(path) {
		app.staticRoutes["POST:"+path] = handler
	} else {
		app.dynamicRoutes = append(app.dynamicRoutes, Route{
			Method:   "POST",
			Pattern:  path,
			Segments: strings.Split(path, "/"),
			Handler:  handler,
		})
	}
}
func PUT[T, V any](app *Application, path string, callback func(Request[T], *Response[V])) {
	handler := &WireHandler[T, V]{
		Path:     path,
		Callback: callback,
	}
	if utils.StaticPath(path) {
		app.staticRoutes["POST:"+path] = handler
	} else {
		app.dynamicRoutes = append(app.dynamicRoutes, Route{
			Method:   "POST",
			Pattern:  path,
			Segments: strings.Split(path, "/"),
			Handler:  handler,
		})
	}
}
func DELETE[T, V any](app *Application, path string, callback func(Request[T], *Response[V])) {
	handler := &WireHandler[T, V]{
		Path:     path,
		Callback: callback,
	}

	if utils.StaticPath(path) {
		app.staticRoutes["POST:"+path] = handler
	} else {
		app.dynamicRoutes = append(app.dynamicRoutes, Route{
			Method:   "POST",
			Pattern:  path,
			Segments: strings.Split(path, "/"),
			Handler:  handler,
		})
	}
}

func PATCH[T, V any](app *Application, path string, callback func(Request[T], *Response[V])) {
	handler := &WireHandler[T, V]{
		Path:     path,
		Callback: callback,
	}
	if utils.StaticPath(path) {
		app.staticRoutes["POST:"+path] = handler
	} else {
		app.dynamicRoutes = append(app.dynamicRoutes, Route{
			Method:   "POST",
			Pattern:  path,
			Segments: strings.Split(path, "/"),
			Handler:  handler,
		})
	}
}
