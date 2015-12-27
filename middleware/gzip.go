package middleware

import (
	"github.com/NYTimes/gziphandler"
	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
	"net/http"
	"sync"
)

var GzipHandler mohttp.Handler = mohttp.HandlerFunc(func(c context.Context) {
	h := gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w2 := &addContentTypeWriter{ResponseWriter: w}
		c2 := mohttp.WithResponseWriter(c, w2)

		mohttp.Next(c2)
	}))

	h.ServeHTTP(mohttp.GetResponseWriter(c), mohttp.GetRequest(c))
})

type addContentTypeWriter struct {
	o sync.Once
	http.ResponseWriter
}

func (w *addContentTypeWriter) Write(b []byte) (int, error) {
	w.o.Do(func() {
		h := w.ResponseWriter.Header()

		if h.Get("Content-Type") == "" {
			h.Set("Content-Type", http.DetectContentType(b))
		}

		h.Del("Content-Length")
	})

	return w.ResponseWriter.Write(b)
}
