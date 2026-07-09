package wire

type Request[T any] struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    T
}
