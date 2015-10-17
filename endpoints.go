package http

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Params  httprouter.Params
	Context context.Context
	Next    Handler
}

type Handler interface {
	Handle(*Context)
}

type HandlerFunc func(*Context)

func (h HandlerFunc) Handle(c *Context) {
	h(c)
}

type Endpoint struct {
	Method   string
	Path     string
	Handlers []Handler
}

func NewEndpoint(method, path string, handlers ...Handler) *Endpoint {
	return &Endpoint{
		Method:   method,
		Path:     path,
		Handlers: handlers,
	}
}

func OPTIONS(path string, handlers ...Handler) *Endpoint {
	return NewEndpoint("OPTIONS", path, handlers...)
}

func GET(path string, handlers ...Handler) *Endpoint {
	return NewEndpoint("GET", path, handlers...)
}

func HEAD(path string, handlers ...Handler) *Endpoint {
	return NewEndpoint("HEAD", path, handlers...)
}

func POST(path string, handlers ...Handler) *Endpoint {
	return NewEndpoint("POST", path, handlers...)
}

func PUT(path string, handlers ...Handler) *Endpoint {
	return NewEndpoint("PUT", path, handlers...)
}

func DELETE(path string, handlers ...Handler) *Endpoint {
	return NewEndpoint("DELETE", path, handlers...)
}
