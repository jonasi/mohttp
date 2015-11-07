package middleware

// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html#sec13.3.3
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html

import (
	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
	"sync"
)

type ETagSource struct {
	etag string
	body []byte
	mu   sync.RWMutex
}

func (e *ETagSource) ETag() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.etag
}

func (e *ETagSource) Body() []byte {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.body
}

func (e *ETagSource) Update(body []byte, etag string) {
	e.mu.Lock()
	defer e.mu.RUnlock()

	e.etag = etag
	e.body = body
}

type ETagDataHandlerFunc func(string, context.Context) (bool, string, interface{}, error)

func (fn ETagDataHandlerFunc) Handle(c context.Context) {
	var (
		h                 = mohttp.GetRequest(c).Header.Get("If-None-Match")
		rw                = mohttp.GetResponseWriter(c)
		match, etag, _, _ = fn(h, c)
	)

	rw.Header().Set("ETag", etag)

	if match {
		rw.WriteHeader(304)
		return
	}
}

func ETagHandler(src *ETagSource) mohttp.Handler {
	return mohttp.HandlerFunc(func(c context.Context) {
		var (
			h  = mohttp.GetRequest(c).Header.Get("If-None-Match")
			rw = mohttp.GetResponseWriter(c)
		)

		if h == "" {
			mohttp.Next(c)
			return
		}

		et := src.ETag()
		rw.Header().Set("ETag", et)

		if h == src.ETag() {
			rw.WriteHeader(304)
			return
		}

		rw.Write(src.body)
	})
}
