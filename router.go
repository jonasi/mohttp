package mohttp

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
)

var notFoundHandler = HandlerFunc(func(c context.Context) {
	GetResponseWriter(c).WriteHeader(http.StatusNotFound)
})

var methodNotAllowedHandler = HandlerFunc(func(c context.Context) {
	GetResponseWriter(c).WriteHeader(http.StatusMethodNotAllowed)
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

func (r *Router) Handle(method, path string, handlers ...Handler) {
	r.router.Handle(method, path, func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		h := append(append([]Handler{}, r.handlers...), handlers...)
		handle(w, req, p, h...)
	})
}

func (r *Router) Use(h ...Handler) {
	r.handlers = append(r.handlers, h...)
}

func (r *Router) Register(routes ...Route) {
	for i := range routes {
		rt := routes[i]

		for _, method := range rt.Methods() {
			for _, path := range rt.Paths() {
				r.Handle(method, path, rt.Handlers()...)
			}
		}
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
	next := HandlerFunc(func(c context.Context) {
		if len(handlers) == 0 {
			return
		}

		cur := handlers[0]
		handlers = handlers[1:]

		cur.Handle(c)
	})

	c := context.Background()
	c = WithRequest(c, req)
	c = WithResponseWriter(c, w)
	c = WithNext(c, next)
	c = WithPathValues(c, &PathValues{
		Params: params(p),
		Query:  query(req.URL.Query()),
	})

	next.Handle(c)
}
