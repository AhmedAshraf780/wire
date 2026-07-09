package wire

import (
	"bufio"
	"fmt"
	"httpserver/internals/utils"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
)

type Application struct {
	Handlers map[string]Handler
}

func NewApplication() *Application {
	RequestsQueue = make(chan RequestMessage, 100)
	ResponseQueue = make(chan ResponseMessage, 100)
	return &Application{
		Handlers: make(map[string]Handler),
	}
}

func (app *Application) Listen(port int) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// runs the workers
	ReqWorker := NewRequestWorker(app)
	ResWorker := NewResponseWorker(app)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go ReqWorker.Work(wg)
	wg.Add(1)
	go ResWorker.Work(wg)
	defer wg.Wait()
	defer listener.Close()

	log.Println("Listening on " + addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go app.handleConnection(conn)
	}
}

func (app *Application) handleConnection(conn net.Conn) {
	log.Println("Handling connection from ", conn.RemoteAddr())
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		// request line
		request := Request[[]byte]{
			Headers: make(map[string]string),
		}

		line, err := reader.ReadString('\n')
		if err == io.EOF {
			log.Println(err)
			return
		}
		if err != nil {
			log.Println(err)
			return
		}
		tokens, err := utils.ParseRequestLine(line)
		if err != nil {
			log.Println(err)
			// TODO: send http response
			return
		}

		request.Method = tokens[0]
		request.Path = tokens[1]
		request.Version = tokens[2]

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
		RequestsQueue <- RequestMessage{
			Request: request,
			Conn:    conn,
		}
		log.Println("Received request from ", conn.RemoteAddr())
	}
}

func GET[T any, V any](app *Application, path string, callback func(Request[T], *Response[V])) {
	handler := &WireHandler[T, V]{
		path,
		callback,
	}
	app.Handlers[utils.GenerateHandlerKey("GET", path)] = handler
}

func POST[T any, V any](app *Application, path string, callback func(Request[T], *Response[V])) {
	handler := &WireHandler[T, V]{
		path,
		callback,
	}
	app.Handlers[utils.GenerateHandlerKey("POST", path)] = handler
}

func PUT[T any, V any](app *Application, path string, callback func(Request[T], *Response[V])) {
	handler := &WireHandler[T, V]{
		path,
		callback,
	}
	app.Handlers[utils.GenerateHandlerKey("PUT", path)] = handler
}
func DELETE[T any, V any](app *Application, path string, callback func(Request[T], *Response[V])) {
	handler := &WireHandler[T, V]{
		path,
		callback,
	}
	app.Handlers[utils.GenerateHandlerKey("DELETE", path)] = handler
}
