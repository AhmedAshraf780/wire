package wire

type Response[T any] struct {
	StatusCode int
	Body       T
	Version    string
	/* headers now don't support duplicate headers like
	/* set-Cookie: a=1
	/* set-Cookie: b=2
	*/
	Headers map[string]string
}

func (r *Response[T]) Header(key string) string {
	return r.Headers[key]
}
func (r *Response[T]) SetHeader(key string, value string) {
	r.Headers[key] = value
}
func (r *Response[T]) Write(statusCode int, body T) error {
	r.StatusCode = statusCode
	r.Body = body
	return nil
}
