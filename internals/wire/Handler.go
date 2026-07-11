package wire

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"

	"github.com/AhmedAshraf780/wire/internals/utils"
)

type handler interface {
	Handle(raw *Request[[]byte], resp *Response[any], conn net.Conn) error
}

type wireHandler[TReq, TRes any] struct {
	Path     string
	Callback func(*Request[TReq], *Response[TRes]) error
}

func (w *wireHandler[TReq, TRes]) Handle(raw *Request[[]byte], resp *Response[any], conn net.Conn) error {
	// check if connection still exists
	value, exist := raw.Headers["Connection"]
	if exist && (value == "Close" || value == "close") {
		println("Close connection")
		resp.SetHeader("Connection", "close")
		conn.Close()
		return nil
	}

	var body TReq

	_, ok := any(body).(EmptyBody)
	if !ok {
		if err := json.Unmarshal(raw.Body, &body); err != nil {
			println("FAILED TO MARSHALING: ", err)
			err := utils.WriteResponse(conn,
				http.StatusBadRequest, "Invalid json body", []byte("Invalid json body"), map[string]string{}, raw.Version,
			)
			if err != nil {
				println("FAILED TO WRITE RESPONSE: ", err)
			}
			return err
		}
	}

	req := Request[TReq]{
		Path:    raw.Path,
		Version: raw.Version,
		Headers: raw.Headers,
		Method:  raw.Method,
		Params:  raw.Params,
		Query:   raw.Query,
		Body:    body,
		Context: raw.Context,
	}

	TypedResp := Response[TRes]{
		Headers: resp.Headers,
	}
	err := w.Callback(&req, &TypedResp)
	if err != nil && errors.Is(err, ErrNext) {
		raw.Context = req.Context
		return ErrNext
	}

	b, err := json.Marshal(TypedResp.Body)
	if err != nil {
		println("FAILED TO MARSHALING: ", err)
		err := utils.WriteResponse(conn,
			http.StatusInternalServerError, "Something Went Wrong", []byte("Something Went Wrong"), map[string]string{}, raw.Version)
		if err != nil {
			println("FAILED TO WRITE RESPONSE: ", err)
		}
		return err
	}
	err = utils.WriteResponse(conn, TypedResp.StatusCode, "success", b, TypedResp.Headers, raw.Version)
	if err != nil {
		println("FAILED TO WRITE RESPONSE: ", err)
	}
	return nil
}
