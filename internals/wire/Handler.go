package wire

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
)

type handler interface {
	Handle(raw *Request[[]byte], resp *Response[any], client Client) error
}

type wireHandler[TReq, TRes any] struct {
	Path     string
	Callback func(*Request[TReq], *Response[TRes]) error
}

func (w *wireHandler[TReq, TRes]) Handle(raw *Request[[]byte], resp *Response[any], client Client) error {
	defer func(conn net.Conn) {
		if client.closed {
			conn.Close()
		}
	}(client.conn)

	var body TReq
	_, ok := any(body).(EmptyBody)
	if !ok {
		if err := json.Unmarshal(raw.Body, &body); err != nil {
			println("FAILED TO MARSHALING: ", err)
			resp := MakeResponse(http.StatusBadRequest, InvalidJsonBody, []byte(InvalidJsonBody), map[string][]string{}, raw.Version)
			client.conn.Write(resp)
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
		Cookies: raw.Cookies,
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
		resp := MakeResponse(http.StatusBadRequest, InvalidJsonBody, []byte(InvalidJsonBody), map[string][]string{}, raw.Version)
		client.conn.Write(resp)
		return err
	}

	res := MakeResponse(TypedResp.StatusCode, Accepted, b, TypedResp.Headers, raw.Version)
	client.conn.Write(res)
	return nil
}
