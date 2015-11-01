package mohttp

import (
	"golang.org/x/net/context"
	"net/http"
)

type StdlibMiddleware func(http.Handler) http.Handler

func FromStdLib(fn StdlibMiddleware) Handler {
	return HandlerFunc(func(c context.Context) {
		fn(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c = WithRequest(c, r)
			c = WithResponseWriter(c, w)
			GetNext(c).Handle(c)
		})).ServeHTTP(GetResponseWriter(c), GetRequest(c))
	})
}

func FromHTTPHandler(handler http.Handler) Handler {
	return HandlerFunc(func(c context.Context) {
		handler.ServeHTTP(GetResponseWriter(c), GetRequest(c))
		GetNext(c).Handle(c)
	})
}
