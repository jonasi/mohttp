package http

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
)

var notFoundHandler = HandlerFunc(func(c *Context) {
	c.Writer.WriteHeader(http.StatusNotFound)
})

var methodNotAllowedHandler = HandlerFunc(func(c *Context) {
	c.Writer.WriteHeader(http.StatusMethodNotAllowed)
})

func NewRouter() *Router {
	r := &Router{
		router:   httprouter.New(),
		handlers: []Handler{},
	}

	r.HandleNotFound(notFoundHandler)
	r.HandleMethodNotAllowed(methodNotAllowedHandler)

	return r
}

type Router struct {
	router   *httprouter.Router
	handlers []Handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) AddGlobalHandler(h ...Handler) {
	r.handlers = append(r.handlers, h...)
}

func (r *Router) Register(endpoints ...*Endpoint) {
	for i := range endpoints {
		ep := endpoints[i]

		r.router.Handle(ep.Method, ep.Path, func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
			h := append(append([]Handler{}, r.handlers...), ep.Handlers...)
			handle(w, req, p, h...)
		})
	}
}

func (r *Router) RegisterHTTPHandler(method, path string, h http.Handler) {
	r.router.Handle(method, path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h.ServeHTTP(w, r)
	})
}

func (r *Router) HandleNotFound(h ...Handler) {
	r.router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handlers := append(append([]Handler{}, r.handlers...), h...)
		handle(w, req, nil, handlers...)
	})
}

func (r *Router) HandleMethodNotAllowed(h ...Handler) {
	r.router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handlers := append(append([]Handler{}, r.handlers...), h...)
		handle(w, req, nil, handlers...)
	})
}

func handle(w http.ResponseWriter, req *http.Request, p httprouter.Params, handlers ...Handler) {
	next := HandlerFunc(func(c *Context) {
		if len(handlers) == 0 {
			return
		}

		cur := handlers[0]
		handlers = handlers[1:]

		cur.Handle(c)
	})

	c := &Context{
		Writer:  w,
		Request: req,
		Params:  p,
		Context: context.Background(),
		Next:    next,
	}

	next.Handle(c)
}
