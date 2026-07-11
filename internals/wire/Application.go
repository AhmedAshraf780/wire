package wire

import (
	"errors"
	"strings"

	"github.com/AhmedAshraf780/wire/internals/utils"
)

type Route struct {
	Method   string
	Pattern  string
	Segments []string
	Handler  Handler
}
type Application struct {
	staticRoutes      map[string]Handler
	dynamicRoutes     []Route
	middlewares       map[string][]Handler
	globalMiddlewares []Handler
}

var WireNext = "next"
var ErrNext = errors.New("next")

func NewApplication() *Application {
	return &Application{
		staticRoutes:      make(map[string]Handler),
		dynamicRoutes:     make([]Route, 0),
		middlewares:       make(map[string][]Handler),
		globalMiddlewares: make([]Handler, 0),
	}
}

func USE[T, V any](app *Application, callback func(*Request[T], *Response[V]) error) {
	app.globalMiddlewares = append(app.globalMiddlewares, &WireHandler[T, V]{
		Path:     "*",
		Callback: callback,
	})
}
func GET[T, V any](app *Application, path string, callbacks ...func(*Request[T], *Response[V]) error) {
	for _, callback := range callbacks {
		handler := &WireHandler[T, V]{
			Path:     path,
			Callback: callback,
		}
		app.middlewares[path] = append(app.middlewares[path], handler)
	}

	handler := &WireHandler[T, V]{
		Path:     path,
		Callback: callbacks[len(callbacks)-1],
	}
	if utils.StaticPath(path) {
		app.staticRoutes["GET:"+path] = handler
	} else {
		app.dynamicRoutes = append(app.dynamicRoutes, Route{
			Method:   "GET",
			Pattern:  path,
			Segments: strings.Split(path, "/"),
			Handler:  handler,
		})
	}
}

func POST[T, V any](app *Application, path string, callbacks ...func(*Request[T], *Response[V]) error) {
	for _, callback := range callbacks {
		handler := &WireHandler[T, V]{
			Path:     path,
			Callback: callback,
		}
		app.middlewares[path] = append(app.middlewares[path], handler)
	}

	handler := &WireHandler[T, V]{
		Path:     path,
		Callback: callbacks[len(callbacks)-1],
	}
	if utils.StaticPath(path) {
		app.staticRoutes["POST:"+path] = handler
	} else {
		app.dynamicRoutes = append(app.dynamicRoutes, Route{
			Method:   "POST",
			Pattern:  path,
			Segments: strings.Split(path, "/"),
			Handler:  handler,
		})
	}
}
func PUT[T, V any](app *Application, path string, callbacks ...func(*Request[T], *Response[V]) error) {
	for _, callback := range callbacks {
		handler := &WireHandler[T, V]{
			Path:     path,
			Callback: callback,
		}
		app.middlewares[path] = append(app.middlewares[path], handler)
	}

	handler := &WireHandler[T, V]{
		Path:     path,
		Callback: callbacks[len(callbacks)-1],
	}
	if utils.StaticPath(path) {
		app.staticRoutes["PUT:"+path] = handler
	} else {
		app.dynamicRoutes = append(app.dynamicRoutes, Route{
			Method:   "PUT",
			Pattern:  path,
			Segments: strings.Split(path, "/"),
			Handler:  handler,
		})
	}
}

func PATCH[T, V any](app *Application, path string, callbacks ...func(*Request[T], *Response[V]) error) {
	for _, callback := range callbacks {
		handler := &WireHandler[T, V]{
			Path:     path,
			Callback: callback,
		}
		app.middlewares[path] = append(app.middlewares[path], handler)
	}

	handler := &WireHandler[T, V]{
		Path:     path,
		Callback: callbacks[len(callbacks)-1],
	}
	if utils.StaticPath(path) {
		app.staticRoutes["PATCH:"+path] = handler
	} else {
		app.dynamicRoutes = append(app.dynamicRoutes, Route{
			Method:   "PATCH",
			Pattern:  path,
			Segments: strings.Split(path, "/"),
			Handler:  handler,
		})
	}
}

func DELETE[T, V any](app *Application, path string, callbacks ...func(*Request[T], *Response[V]) error) {
	for _, callback := range callbacks {
		handler := &WireHandler[T, V]{
			Path:     path,
			Callback: callback,
		}
		app.middlewares[path] = append(app.middlewares[path], handler)
	}

	handler := &WireHandler[T, V]{
		Path:     path,
		Callback: callbacks[len(callbacks)-1],
	}
	if utils.StaticPath(path) {
		app.staticRoutes["DELETE:"+path] = handler
	} else {
		app.dynamicRoutes = append(app.dynamicRoutes, Route{
			Method:   "DELETE",
			Pattern:  path,
			Segments: strings.Split(path, "/"),
			Handler:  handler,
		})
	}
}
