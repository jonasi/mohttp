package http

import (
	"net/http"
)

func StripPrefix(prefix string) Handler {
	return FromStdLib(func(h http.Handler) http.Handler {
		return http.StripPrefix(prefix, h)
	})
}

func Prefix(prefix string, endpoints ...Endpoint) []Endpoint {
	eps := make([]Endpoint, len(endpoints))

	for i, ep := range endpoints {
		p := ep.Paths()
		p2 := make([]string, len(p))

		for i := range p {
			p2[i] = prefix + p[i]
		}

		eps[i] = NewComplexEndpoint(ep.Methods(), p2, ep.Handlers()...)
	}

	return eps
}
