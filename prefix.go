package mohttp

import (
	"net/http"
)

func StripPrefix(prefix string) Handler {
	return FromStdLib(func(h http.Handler) http.Handler {
		return http.StripPrefix(prefix, h)
	})
}

func Prefix(prefix string, routes ...Route) []Route {
	rts := make([]Route, len(routes))

	for i, rt := range routes {
		p := rt.Paths()
		p2 := make([]string, len(p))

		for i := range p {
			p2[i] = prefix + p[i]
		}

		rts[i] = NewComplexRoute(rt.Methods(), p2, rt.Handlers()...)
	}

	return rts
}
