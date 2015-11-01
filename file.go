package mohttp

import (
	"golang.org/x/net/context"
	"net/http"
)

func FileHandler(fn func(context.Context) string) Handler {
	return HandlerFunc(func(c context.Context) {
		path := fn(c)
		http.ServeFile(GetResponseWriter(c), GetRequest(c), path)
	})
}
