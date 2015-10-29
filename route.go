package mohttp

type Handler interface {
	Handle(*Context)
}

type HandlerFunc func(*Context)

func (h HandlerFunc) Handle(c *Context) {
	h(c)
}

type Route interface {
	Methods() []string
	Paths() []string
	Handlers() []Handler
}

type route struct {
	methods  []string
	paths    []string
	handlers []Handler
}

func (e *route) Methods() []string   { return e.methods }
func (e *route) Paths() []string     { return e.paths }
func (e *route) Handlers() []Handler { return e.handlers }

func NewComplexRoute(methods []string, paths []string, handlers ...Handler) Route {
	return &route{
		methods:  methods,
		paths:    paths,
		handlers: handlers,
	}
}

func NewRoute(method, path string, handlers ...Handler) Route {
	return &route{
		methods:  []string{method},
		paths:    []string{path},
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

func DELETE(path string, handlers ...Handler) Route {
	return NewRoute("DELETE", path, handlers...)
}

func ALL(path string, handlers ...Handler) Route {
	return &route{
		methods:  []string{"OPTIONS", "GET", "HEAD", "POST", "PUT", "DELETE"},
		paths:    []string{path},
		handlers: handlers,
	}
}
