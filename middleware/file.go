package middleware

import (
	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
	"net/http"
)

func FileHandler(fn func(context.Context) string) mohttp.Handler {
	return mohttp.HandlerFunc(func(c context.Context) {
		path := fn(c)
		http.ServeFile(mohttp.GetResponseWriter(c), mohttp.GetRequest(c), path)
	})
}
