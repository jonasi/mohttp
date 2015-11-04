package mohttp

import (
	"net/http"
)

func StripPrefixHandler(prefix string) Handler {
	return FromHTTPMiddleware(func(h http.Handler) http.Handler {
		return http.StripPrefix(prefix, h)
	})
}

func Prefix(prefix string, routes ...Route) []Route {
	rts := make([]Route, len(routes))

	for i, rt := range routes {
		rts[i] = NewRoute(rt.Method(), prefix+rt.Path(), rt.Handlers()...)
	}

	return rts
}
