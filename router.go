package mohttp

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

var notFoundHandler = HandlerFunc(func(c context.Context) {
	GetResponseWriter(c).WriteHeader(http.StatusNotFound)
})

var methodNotAllowedHandler = HandlerFunc(func(c context.Context) {
	GetResponseWriter(c).WriteHeader(http.StatusMethodNotAllowed)
})

func NewRouter() *Router {
	r := &Router{
		router: httprouter.New(),
		use:    []Handler{},
	}

	r.HandleNotFound(notFoundHandler)
	r.HandleMethodNotAllowed(methodNotAllowedHandler)

	return r
}

type Router struct {
	router *httprouter.Router
	use    []Handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) Handle(method, path string, handlers ...Handler) {
	r.router.Handle(method, path, func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		h := append(append([]Handler{}, r.use...), handlers...)
		c := HTTPContext(w, req, p)
		Serve(c, h...)
	})
}

func (r *Router) Use(h ...Handler) {
	r.use = append(r.use, h...)
}

func (r *Router) Register(routes ...Route) {
	for _, rt := range routes {
		r.Handle(rt.Method(), rt.Path(), rt.Handlers()...)
	}
}

func (r *Router) RegisterHTTPHandler(method, path string, h http.Handler) {
	r.router.Handle(method, path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h.ServeHTTP(w, r)
	})
}

func (r *Router) HandleNotFound(h ...Handler) {
	r.router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handlers := append(append([]Handler{}, r.use...), h...)
		c := HTTPContext(w, req, nil)
		Serve(c, handlers...)
	})
}

func (r *Router) HandleMethodNotAllowed(h ...Handler) {
	r.router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handlers := append(append([]Handler{}, r.use...), h...)
		c := HTTPContext(w, req, nil)
		Serve(c, handlers...)
	})
}
