package mohttp

import (
	"golang.org/x/net/context"
	"net/http"
)

func Redirect(path string) Handler {
	return HandlerFunc(func(c context.Context) {
		http.Redirect(GetResponseWriter(c), GetRequest(c), path, http.StatusTemporaryRedirect)
		GetNext(c).Handle(c)
	})
}
