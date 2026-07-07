package wire

import (
	"httpserver/internals/utils"
	"log"
	"sync"
)

type ResponseWorker struct {
	App *Application
}

func NewResponseWorker(app *Application) *ResponseWorker {
	return &ResponseWorker{app}
}

func (rw *ResponseWorker) Work(wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range ResponseQueue {
		log.Println("Response queue is working ")
		res := utils.MakeResponse(t.Response.StatusCode, "Success", t.Response.Body, t.Response.Headers, t.Response.Version)
		log.Println(string(res))
		_, err := t.Conn.Write(res)
		if err != nil {
			log.Println(err)
		}
	}
}
