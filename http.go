package mohttp

import (
	"golang.org/x/net/context"
	"net/http"
)

func TemporaryRedirect(c context.Context, path string) {
	http.Redirect(GetResponseWriter(c), GetRequest(c), path, http.StatusTemporaryRedirect)
}

func TemporaryRedirectHandler(path string) Handler {
	return HandlerFunc(func(c context.Context) {
		TemporaryRedirect(c, path)
		Next(c)
	})
}

func PermanentRedirect(c context.Context, path string) {
	http.Redirect(GetResponseWriter(c), GetRequest(c), path, http.StatusMovedPermanently)
}

func PermanentRedirectHandler(path string) Handler {
	return HandlerFunc(func(c context.Context) {
		PermanentRedirect(c, path)
		Next(c)
	})
}

func Status(c context.Context, code int) {
	GetResponseWriter(c).WriteHeader(code)
}

func StatusHandler(code int) Handler {
	return HandlerFunc(func(c context.Context) {
		Status(c, code)
		Next(c)
	})
}

func HeadersHandler(pairs ...string) Handler {
	l := len(pairs)
	if l%2 == 1 {
		panic("Header pairs must be a multiple of 2")
	}

	return HandlerFunc(func(c context.Context) {
		for i := 0; i < l/2; i++ {
			GetResponseWriter(c).Header().Add(pairs[2*i], pairs[2*i+1])
		}

		Next(c)
	})
}
