package middleware

// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html#sec13.3.3
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html

import (
	"bytes"
	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"time"
)

type timedCacheResp struct {
	lastSet time.Time
	code    int
	headers http.Header
	body    *bytes.Reader
	sync.Mutex
}

type TimedCache struct {
	Duration time.Duration
	caches   map[string]*timedCacheResp
	mu       sync.Mutex
}

func (t *TimedCache) Handle(c context.Context) {
	id := mohttp.RouteID(c) + "." + queryID(c)

	t.mu.Lock()
	if t.caches == nil {
		t.caches = map[string]*timedCacheResp{}
	}

	resp, ok := t.caches[id]
	if !ok {
		t.caches[id] = &timedCacheResp{}
		resp = t.caches[id]
	}

	t.mu.Unlock()

	resp.Lock()
	defer resp.Unlock()

	rw := mohttp.GetResponseWriter(c)

	if time.Since(resp.lastSet) > t.Duration {
		rec := httptest.NewRecorder()
		c = mohttp.WithResponseWriter(c, rec)
		mohttp.Next(c)

		resp.lastSet = time.Now()
		resp.code = rec.Code
		resp.headers = rec.HeaderMap
		resp.body = bytes.NewReader(rec.Body.Bytes())
	}

	copyResp(resp.code, resp.headers, resp.body, rw)
}

func copyResp(code int, headers http.Header, body *bytes.Reader, newRw http.ResponseWriter) {
	h := newRw.Header()

	for k, v := range headers {
		h[k] = append(h[k], v...)
	}

	newRw.WriteHeader(code)
	body.WriteTo(newRw)
	body.Seek(0, 0)
}

func queryID(c context.Context) string {
	q := mohttp.GetRequest(c).URL.Query()
	vals := make([]string, len(q))
	i := 0

	for k, v := range q {
		sort.Strings(v)
		vals[i] = k + "=" + strings.Join(v, "&")
		i++
	}

	sort.Strings(vals)

	return strings.Join(vals, ",")
}

func EtagHandlerFunc(fn func(context.Context) (interface{}, string, error)) mohttp.Handler {
	return mohttp.DataHandlerFunc(func(c context.Context) (interface{}, error) {
		d, etag, err := fn(c)

		if err != nil {
			return nil, err
		}

		var (
			found = mohttp.GetRequest(c).Header.Get("If-None-Match")
			rw    = mohttp.GetResponseWriter(c)
		)

		rw.Header().Set("ETag", etag)

		if etag == found {
			rw.WriteHeader(304)
			return mohttp.DataNoBody, nil
		}

		return d, nil
	})
}
