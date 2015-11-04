package mohttp

import (
	"golang.org/x/net/context"
)

type PhaseHandler interface {
	BeforeMain(context.Context)
	Handler
}

type Handler interface {
	Handle(context.Context)
}

type HandlerFunc func(context.Context)

func (h HandlerFunc) Handle(c context.Context) {
	h(c)
}

type BeforeHandlerFunc func(context.Context)

func (h BeforeHandlerFunc) BeforeMain(c context.Context) {
	h(c)
}

func (h BeforeHandlerFunc) Handle(c context.Context) {
	Next(c)
}

type Route interface {
	Path() string
	Method() string
	Handlers() []Handler
}

type route struct {
	method   string
	path     string
	handlers []Handler
}

func (e *route) Method() string      { return e.method }
func (e *route) Path() string        { return e.path }
func (e *route) Handlers() []Handler { return e.handlers }

func NewRoute(method, path string, handlers ...Handler) Route {
	return &route{
		method:   method,
		path:     path,
		handlers: handlers,
	}
}

func OPTIONS(path string, handlers ...Handler) Route {
	return NewRoute("OPTIONS", path, handlers...)
}

func GET(path string, handlers ...Handler) Route {
	return NewRoute("GET", path, handlers...)
}

func HEAD(path string, handlers ...Handler) Route {
	return NewRoute("HEAD", path, handlers...)
}

func POST(path string, handlers ...Handler) Route {
	return NewRoute("POST", path, handlers...)
}

func PUT(path string, handlers ...Handler) Route {
	return NewRoute("PUT", path, handlers...)
}

func PATCH(path string, handlers ...Handler) Route {
	return NewRoute("PATCH", path, handlers...)
}

func DELETE(path string, handlers ...Handler) Route {
	return NewRoute("DELETE", path, handlers...)
}

func TRACE(path string, handlers ...Handler) Route {
	return NewRoute("TRACE", path, handlers...)
}

func CONNECT(path string, handlers ...Handler) Route {
	return NewRoute("CONNECT", path, handlers...)
}

var methods = [...]string{"OPTIONS", "GET", "HEAD", "POST", "PUT", "DELETE", "PATCH", "TRACE", "CONNECT"}

func ALL(path string, handlers ...Handler) []Route {
	rts := make([]Route, len(methods))

	for i := range methods {
		rts[i] = NewRoute(methods[i], path, handlers...)
	}

	return rts
}
