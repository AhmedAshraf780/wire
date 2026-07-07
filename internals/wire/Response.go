package wire

type Response struct {
	StatusCode int
	Body       []byte
	Version    string
	Headers    map[string]string
}
