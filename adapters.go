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
		h := fn(mkNextHTTPHandler(c))
		h.ServeHTTP(GetResponseWriter(c), GetRequest(c))
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

type NegronHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, http.HandlerFunc)
}

type NegroniHandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (fn NegroniHandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fn(rw, r, next)
}

func FromNegroniHandler(n NegronHandler) Handler {
	return HandlerFunc(func(c context.Context) {
		n.ServeHTTP(GetResponseWriter(c), GetRequest(c), mkNextHTTPHandler(c))
	})
}

func FromNegroniHandlerFunc(fn func(http.ResponseWriter, *http.Request, http.HandlerFunc)) Handler {
	return FromNegroniHandler(NegroniHandlerFunc(fn))
}

func mkNextHTTPHandler(c context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c = WithRequest(c, r)
		c = WithResponseWriter(c, w)
		Next(c)
	}
}
