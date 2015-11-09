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

func CheckETag(c context.Context, etag string) bool {
	if etag == mohttp.GetRequest(c).Header.Get("If-None-Match") {
		mohttp.GetResponseWriter(c).WriteHeader(304)
		return true
	}

	return false
}

type ETagCache struct {
	etag string
	body []byte
	mu   sync.RWMutex
}

func (e *ETagCache) ETag() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.etag
}

func (e *ETagCache) Body() []byte {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.body
}

func (e *ETagCache) Update(body []byte, etag string) {
	e.mu.Lock()
	defer e.mu.RUnlock()

	e.etag = etag
	e.body = body
}

func (e *ETagCache) Handle(c context.Context) {
	var (
		h  = mohttp.GetRequest(c).Header.Get("If-None-Match")
		rw = mohttp.GetResponseWriter(c)
		et = e.ETag()
	)

	if et != "" {
		rw.Header().Set("ETag", et)
	}

	if h == "" {
		mohttp.Next(c)
		return
	}

	if h == e.ETag() {
		rw.WriteHeader(304)
		return
	}

	rw.Write(e.Body())
}
