package wire

import (
	"encoding/json"
	"strconv"
)

type Handler interface {
	Handle(raw Request[[]byte], resp *Response[any]) (*Response[[]byte], error)
}

type WireHandler[TReq, TRes any] struct {
	Path     string
	Callback func(Request[TReq], *Response[TRes])
}

func (w *WireHandler[TReq, TRes]) Handle(raw Request[[]byte], resp *Response[any]) (*Response[[]byte], error) {
	var body TReq

	if err := json.Unmarshal(raw.Body, &body); err != nil {
		return &Response[[]byte]{
			Version:    raw.Version,
			StatusCode: 400,
			Headers:    map[string]string{},
			Body:       []byte(err.Error()),
		}, err
	}

	req := Request[TReq]{
		Path:    raw.Path,
		Version: raw.Version,
		Headers: raw.Headers,
		Method:  raw.Method,
		Body:    body,
	}

	TypedResp := Response[TRes]{
		Headers: resp.Headers,
	}

	w.Callback(req, &TypedResp)

	b, err := json.Marshal(TypedResp.Body)
	if err != nil {
		return nil, err
	}
	TypedResp.Headers["Content-Length"] = strconv.Itoa(len(b))
	TypedResp.Headers["Content-type"] = "application/json"

	return &Response[[]byte]{
		Version:    raw.Version,
		Headers:    TypedResp.Headers,
		StatusCode: TypedResp.StatusCode,
		Body:       b,
	}, nil
}
