package wire

import (
	"errors"
	"strings"
)

type route struct {
	Method   string
	Pattern  string
	Segments []string
	Handler  handler
}
type Application struct {
	staticRoutes      map[string]handler
	dynamicRoutes     []route
	middlewares       map[string][]handler
	globalMiddlewares []handler
}

var ErrNext = errors.New("next")

func NewApplication() *Application {
	return &Application{
		staticRoutes:      make(map[string]handler),
		dynamicRoutes:     make([]route, 0),
		middlewares:       make(map[string][]handler),
		globalMiddlewares: make([]handler, 0),
	}
}

func USE[T, V any](app *Application, callback func(*Request[T], *Response[V]) error) {
	app.globalMiddlewares = append(app.globalMiddlewares, &wireHandler[T, V]{
		Path:     "*",
		Callback: callback,
	})
}
func GET[T, V any](app *Application, path string, callbacks ...func(*Request[T], *Response[V]) error) {
	for _, callback := range callbacks {
		handler := &wireHandler[T, V]{
			Path:     path,
			Callback: callback,
		}
		app.middlewares[generateHandlerKey("GET", path)] = append(app.middlewares[path], handler)
	}

	handler := &wireHandler[T, V]{
		Path:     path,
		Callback: callbacks[len(callbacks)-1],
	}
	if staticPath(path) {
		app.staticRoutes["GET:"+path] = handler
	} else {
		app.dynamicRoutes = append(app.dynamicRoutes, route{
			Method:   "GET",
			Pattern:  path,
			Segments: strings.Split(path, "/"),
			Handler:  handler,
		})
	}
}

func POST[T, V any](app *Application, path string, callbacks ...func(*Request[T], *Response[V]) error) {
	for _, callback := range callbacks {
		handler := &wireHandler[T, V]{
			Path:     path,
			Callback: callback,
		}
		app.middlewares[generateHandlerKey("POST", path)] = append(app.middlewares[path], handler)
	}

	handler := &wireHandler[T, V]{
		Path:     path,
		Callback: callbacks[len(callbacks)-1],
	}
	if staticPath(path) {
		app.staticRoutes["POST:"+path] = handler
	} else {
		app.dynamicRoutes = append(app.dynamicRoutes, route{
			Method:   "POST",
			Pattern:  path,
			Segments: strings.Split(path, "/"),
			Handler:  handler,
		})
	}
}
func PUT[T, V any](app *Application, path string, callbacks ...func(*Request[T], *Response[V]) error) {
	for _, callback := range callbacks {
		handler := &wireHandler[T, V]{
			Path:     path,
			Callback: callback,
		}
		app.middlewares[generateHandlerKey("PUT", path)] = append(app.middlewares[path], handler)
	}

	handler := &wireHandler[T, V]{
		Path:     path,
		Callback: callbacks[len(callbacks)-1],
	}
	if staticPath(path) {
		app.staticRoutes["PUT:"+path] = handler
	} else {
		app.dynamicRoutes = append(app.dynamicRoutes, route{
			Method:   "PUT",
			Pattern:  path,
			Segments: strings.Split(path, "/"),
			Handler:  handler,
		})
	}
}

func PATCH[T, V any](app *Application, path string, callbacks ...func(*Request[T], *Response[V]) error) {
	for _, callback := range callbacks {
		handler := &wireHandler[T, V]{
			Path:     path,
			Callback: callback,
		}
		app.middlewares[generateHandlerKey("PATCH", path)] = append(app.middlewares[path], handler)
	}

	handler := &wireHandler[T, V]{
		Path:     path,
		Callback: callbacks[len(callbacks)-1],
	}
	if staticPath(path) {
		app.staticRoutes["PATCH:"+path] = handler
	} else {
		app.dynamicRoutes = append(app.dynamicRoutes, route{
			Method:   "PATCH",
			Pattern:  path,
			Segments: strings.Split(path, "/"),
			Handler:  handler,
		})
	}
}

func DELETE[T, V any](app *Application, path string, callbacks ...func(*Request[T], *Response[V]) error) {
	for _, callback := range callbacks {
		handler := &wireHandler[T, V]{
			Path:     path,
			Callback: callback,
		}
		app.middlewares[generateHandlerKey("DELETE", path)] = append(app.middlewares[path], handler)
	}

	handler := &wireHandler[T, V]{
		Path:     path,
		Callback: callbacks[len(callbacks)-1],
	}
	if staticPath(path) {
		app.staticRoutes["DELETE:"+path] = handler
	} else {
		app.dynamicRoutes = append(app.dynamicRoutes, route{
			Method:   "DELETE",
			Pattern:  path,
			Segments: strings.Split(path, "/"),
			Handler:  handler,
		})
	}
}
