package http

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
