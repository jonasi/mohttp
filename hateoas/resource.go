package hateoas

import (
	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
)

// http://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html
// https://en.wikipedia.org/wiki/Representational_state_transfer
// https://en.wikipedia.org/wiki/HATEOAS

type link struct {
	Rel      string
	Resource *Resource
}

var setResource, resStore = mohttp.NewContextValuePair("github.com/jonasi/mohttp/hateoas.Resource")

func GetResource(c context.Context) *Resource {
	return resStore.Get(c).(*Resource)
}

type ResourceOption func(*Resource)

func NewResource(opts ...ResourceOption) *Resource {
	r := &Resource{
		handlers: map[string][]mohttp.Handler{},
	}

	r.use = []mohttp.Handler{
		setResource(r),
	}

	for i := range opts {
		opts[i](r)
	}

	return r
}

type Resource struct {
	path     string
	use      []mohttp.Handler
	handlers map[string][]mohttp.Handler
	links    []link
}

func (r *Resource) Routes() []mohttp.Route {
	if len(r.handlers["GET"]) > 0 && len(r.handlers["HEAD"]) == 0 {
		r.handlers["HEAD"] = r.handlers["GET"]
	}

	var (
		rts = make([]mohttp.Route, len(r.handlers))
		i   = 0
	)

	for method := range r.handlers {
		h := make([]mohttp.Handler, len(r.use)+len(r.handlers[method]))
		copy(h, r.use)
		copy(h[len(r.use):], r.handlers[method])

		rts[i] = mohttp.NewRoute(method, r.path, h...)
		i++
	}

	return rts
}

func Path(path string) ResourceOption {
	return func(r *Resource) { r.path = path }
}

func Use(h ...mohttp.Handler) ResourceOption {
	return func(r *Resource) { r.use = append(r.use, h...) }
}

func Link(rel string, l *Resource) ResourceOption {
	return func(r *Resource) { r.links = append(r.links, link{rel, l}) }
}

func OPTIONS(h ...mohttp.Handler) ResourceOption {
	return addHandlers("OPTIONS", h...)
}

func HEAD(h ...mohttp.Handler) ResourceOption {
	return addHandlers("HEAD", h...)
}

func GET(h ...mohttp.Handler) ResourceOption {
	return addHandlers("GET", h...)
}

func POST(h ...mohttp.Handler) ResourceOption {
	return addHandlers("POST", h...)
}

func PUT(h ...mohttp.Handler) ResourceOption {
	return addHandlers("PUT", h...)
}

func PATCH(h ...mohttp.Handler) ResourceOption {
	return addHandlers("PATCH", h...)
}

func DELETE(h ...mohttp.Handler) ResourceOption {
	return addHandlers("DELETE", h...)
}

func addHandlers(method string, h ...mohttp.Handler) ResourceOption {
	return func(r *Resource) {
		if _, ok := r.handlers[method]; !ok {
			r.handlers[method] = []mohttp.Handler{}
		}

		r.handlers[method] = append(r.handlers[method], h...)
	}
}
