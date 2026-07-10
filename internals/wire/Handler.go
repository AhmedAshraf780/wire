package wire

import (
	"encoding/json"
	"httpserver/internals/utils"
	"net"
	"net/http"
)

type Handler interface {
	Handle(raw Request[[]byte], resp *Response[any], conn net.Conn) error
}

type WireHandler[TReq, TRes any] struct {
	Path     string
	Callback func(Request[TReq], *Response[TRes])
}

func (w *WireHandler[TReq, TRes]) Handle(raw Request[[]byte], resp *Response[any], conn net.Conn) error {
	var body TReq

	_, ok := any(body).(EmptyBody)
	if !ok {
		if err := json.Unmarshal(raw.Body, &body); err != nil {
			resp := utils.MakeResponse(
				http.StatusBadRequest, "Invalid json body", []byte("Invalid json body"), map[string]string{}, raw.Version,
			)
			conn.Write(resp)
			return err
		}
	}

	// check the mode first
	req := Request[TReq]{
		Path:    raw.Path,
		Version: raw.Version,
		Headers: raw.Headers,
		Method:  raw.Method,
		Body:    body,
	}

	// extract the params from the path
	params, err := utils.ExtractParams(w.Path, req.Path)
	if err != nil {
		resp := utils.MakeResponse(
			http.StatusBadRequest, "Invalid json body", []byte("Invalid request params"), map[string]string{}, raw.Version)
		conn.Write(resp)
		return err
	}
	req.Params = params

	// extract the query
	quries := utils.ParseQuery(req.Path)
	req.Query = quries

	TypedResp := Response[TRes]{
		Headers: resp.Headers,
	}
	w.Callback(req, &TypedResp)

	b, err := json.Marshal(TypedResp.Body)
	if err != nil {
		resp := utils.MakeResponse(
			http.StatusBadRequest, "Invalid json body", []byte("Invalid json body"), map[string]string{}, raw.Version)
		conn.Write(resp)
		return err
	}
	res := utils.MakeResponse(TypedResp.StatusCode, "success", b, TypedResp.Headers, raw.Version)
	conn.Write(res)
	return nil
}
