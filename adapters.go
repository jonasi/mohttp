package mohttp

import (
	"golang.org/x/net/context"
	"net/http"
)

// FromHTTPMiddleware transforms a golang stdlib middleware of the form
//   func(http.Handler) Handler
// into a mohttp.Handler
func FromHTTPMiddleware(fn func(http.Handler) http.Handler) Handler {
	return HandlerFunc(func(c context.Context) {
		fn(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c = WithRequest(c, r)
			c = WithResponseWriter(c, w)
			Next(c)
		})).ServeHTTP(GetResponseWriter(c), GetRequest(c))
	})
}

// FromHTTPHandler transforms a golang stdlib http.Handler
// into a mohttp.Handler
func FromHTTPHandler(handler http.Handler) Handler {
	return HandlerFunc(func(c context.Context) {
		handler.ServeHTTP(GetResponseWriter(c), GetRequest(c))
		Next(c)
	})
}
