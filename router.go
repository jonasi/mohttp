package http

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
)

func NewRouter() *Router {
	return &Router{
		router:   httprouter.New(),
		handlers: []Handler{},
	}
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
	for _, ep := range endpoints {
		r.router.Handle(ep.Method, ep.Path, r.mkHandle(ep.Handlers...))
	}
}

func (r *Router) RegisterHTTPHandler(method, path string, h http.Handler) {
	r.router.Handle(method, path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h.ServeHTTP(w, r)
	})
}

func (r *Router) HandleNotFound(h Handler) {
	r.router.NotFound = r.mkHandler(h)
}

func (r *Router) HandleMethodNotAllowed(h Handler) {
	r.router.MethodNotAllowed = r.mkHandler(h)
}

func (r *Router) mkHandle(h ...Handler) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		r.handle(w, req, p, h...)
	}
}

func (r *Router) mkHandler(h ...Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r.handle(w, req, nil, h...)
	})
}

func (r *Router) handle(w http.ResponseWriter, req *http.Request, p httprouter.Params, handlers ...Handler) {
	handlers = append(r.handlers, handlers...)

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
