package hateoas

import (
	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
)

var setService, svcStore = mohttp.NewContextValuePair("github.com/jonasi/mohttp/hateoas.Service")

type ServiceOption func(*Service)

func AddResource(r ...*Resource) ServiceOption {
	return func(s *Service) {
		s.resources = append(s.resources, r...)
	}
}

func ServiceUse(h ...mohttp.Handler) ServiceOption {
	return func(s *Service) {
		s.Use(h...)
	}
}

func GetService(c context.Context) (*Service, bool) {
	svc, ok := svcStore.Get(c).(*Service)
	return svc, ok
}

func NewService(opts ...ServiceOption) *Service {
	s := &Service{
		resources: []*Resource{},
	}

	s.use = []mohttp.Handler{
		setService(s),
	}

	for i := range opts {
		opts[i](s)
	}

	return s
}

type Service struct {
	resources []*Resource
	use       []mohttp.Handler
}

func (s *Service) Use(h ...mohttp.Handler) {
	s.use = append(s.use, h...)
}

func (s *Service) Routes() []mohttp.Route {
	allRoutes := []mohttp.Route{}

	for _, r := range s.resources {
		rts := r.Routes()

		for i, rt := range rts {
			old := rt.Handlers()
			h := make([]mohttp.Handler, len(s.use)+len(old))
			copy(h, s.use)
			copy(h[len(s.use):], old)

			rts[i] = mohttp.NewRoute(rt.Method(), rt.Path(), h...)
		}

		allRoutes = append(allRoutes, rts...)
	}

	return allRoutes
}
