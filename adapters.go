package mohttp

import (
	"net/http"
)

type StdlibMiddleware func(http.Handler) http.Handler

func FromStdLib(fn StdlibMiddleware) Handler {
	return HandlerFunc(func(c *Context) {
		fn(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Writer = w
			c.Request = r

			c.Next.Handle(c)
		})).ServeHTTP(c.Writer, c.Request)
	})
}

func FromHTTPHandler(handler http.Handler) Handler {
	return HandlerFunc(func(c *Context) {
		handler.ServeHTTP(c.Writer, c.Request)
		c.Next.Handle(c)
	})
}
