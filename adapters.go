package mohttp

import (
	"golang.org/x/net/context"
	"net/http"
)

func FromHTTPMiddleware(fn func(http.Handler) http.Handler) Handler {
	return HandlerFunc(func(c context.Context) {
		fn(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c = WithRequest(c, r)
			c = WithResponseWriter(c, w)
			Next(c)
		})).ServeHTTP(GetResponseWriter(c), GetRequest(c))
	})
}

func FromHTTPHandler(handler http.Handler) Handler {
	return HandlerFunc(func(c context.Context) {
		handler.ServeHTTP(GetResponseWriter(c), GetRequest(c))
		Next(c)
	})
}
