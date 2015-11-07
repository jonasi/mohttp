package middleware

import (
	"github.com/jonasi/mohttp"
	"net/http"
)

func StripPrefixHandler(prefix string) mohttp.Handler {
	return mohttp.FromHTTPMiddleware(func(h http.Handler) http.Handler {
		return http.StripPrefix(prefix, h)
	})
}

func Prefix(prefix string, routes ...mohttp.Route) []mohttp.Route {
	rts := make([]mohttp.Route, len(routes))

	for i, rt := range routes {
		rts[i] = mohttp.NewRoute(rt.Method(), prefix+rt.Path(), rt.Handlers()...)
	}

	return rts
}
