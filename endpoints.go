package http

type Handler interface {
	Handle(*Context)
}

type HandlerFunc func(*Context)

func (h HandlerFunc) Handle(c *Context) {
	h(c)
}

type Endpoint interface {
	Methods() []string
	Paths() []string
	Handlers() []Handler
}

type endpoint struct {
	methods  []string
	paths    []string
	handlers []Handler
}

func (e *endpoint) Methods() []string   { return e.methods }
func (e *endpoint) Paths() []string     { return e.paths }
func (e *endpoint) Handlers() []Handler { return e.handlers }

func NewComplexEndpoint(methods []string, paths []string, handlers ...Handler) Endpoint {
	return &endpoint{
		methods:  methods,
		paths:    paths,
		handlers: handlers,
	}
}

func NewEndpoint(method, path string, handlers ...Handler) Endpoint {
	return &endpoint{
		methods:  []string{method},
		paths:    []string{path},
		handlers: handlers,
	}
}

func OPTIONS(path string, handlers ...Handler) Endpoint {
	return NewEndpoint("OPTIONS", path, handlers...)
}

func GET(path string, handlers ...Handler) Endpoint {
	return NewEndpoint("GET", path, handlers...)
}

func HEAD(path string, handlers ...Handler) Endpoint {
	return NewEndpoint("HEAD", path, handlers...)
}

func POST(path string, handlers ...Handler) Endpoint {
	return NewEndpoint("POST", path, handlers...)
}

func PUT(path string, handlers ...Handler) Endpoint {
	return NewEndpoint("PUT", path, handlers...)
}

func DELETE(path string, handlers ...Handler) Endpoint {
	return NewEndpoint("DELETE", path, handlers...)
}

func ALL(path string, handlers ...Handler) Endpoint {
	return &endpoint{
		methods:  []string{"OPTIONS", "GET", "HEAD", "POST", "PUT", "DELETE"},
		paths:    []string{path},
		handlers: handlers,
	}
}
