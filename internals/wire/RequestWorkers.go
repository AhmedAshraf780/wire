package wire

import (
	"httpserver/internals/utils"
	"log"
	"sync"
)

type RequestWorker struct {
	app *Application
}

func NewRequestWorker(app *Application) *RequestWorker {
	return &RequestWorker{app: app}
}

func (r *RequestWorker) Work(wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range RequestsQueue {
		log.Println("I am working on task", t.Request.Path)
		key := utils.GenerateHandlerKey(t.Request.Method, t.Request.Path)
		handler, ok := r.app.Handlers[key]
		if !ok {
			// publish response error with not found
			log.Println("No handler for task", t.Request.Path)
			// TODO: add mutex here
			ResponseQueue <- ResponseMessage{
				Conn: t.Conn,
				Response: Response{
					// TODO: add standard errors
					StatusCode: 404,
					Body:       []byte{},
					// TODO: Handle headers
					Headers: utils.DefaultHeaders(),
					Version: t.Request.Version,
				},
			}
			continue
		}
		// execute the handler
		resp := &Response{
			Version: t.Request.Version,
			Headers: map[string]string{},
		}
		handler.Callback(t.Request, resp)
		ResponseQueue <- ResponseMessage{
			Conn:     t.Conn,
			Response: *resp,
		}
	}
}
