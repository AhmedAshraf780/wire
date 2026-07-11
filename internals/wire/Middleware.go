package wire

import "errors"

type WireMiddleware struct {
	Path        string
	Handler     Handler
	NextHandler Handler
}

func Next() error {
	return errors.New(WireNext)
}
