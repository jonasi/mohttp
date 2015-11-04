package mohttp

import (
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"time"
)

type RequestSummary struct {
	StartTime     time.Time
	Protocol      string
	Method        string
	URL           *url.URL
	ClientIP      string
	UserAgent     string
	Referer       string
	Duration      time.Duration
	StatusCode    int
	ContentLength int
}

func RequestLogger(fn func(*RequestSummary)) Handler {
	return HandlerFunc(func(c context.Context) {
		req := GetRequest(c)

		st := &RequestSummary{
			StartTime: time.Now(),
			Protocol:  req.Proto,
			Method:    req.Method,
			URL:       req.URL,
			UserAgent: req.UserAgent(),
			ClientIP:  req.RemoteAddr,
			Referer:   req.Referer(),
		}

		rw := newStatsResponseWriter(GetResponseWriter(c))
		c = WithResponseWriter(c, rw)

		Next(c)

		st.Duration = time.Since(st.StartTime)
		st.StatusCode = rw.Status()
		st.ContentLength = rw.ContentLength()

		fn(st)
	})
}

func newStatsResponseWriter(rw http.ResponseWriter) *statsResponseWriter {
	return &statsResponseWriter{
		ResponseWriter: rw,
		status:         200,
	}
}

type statsResponseWriter struct {
	http.ResponseWriter
	status        int
	contentLength int
}

func (rw *statsResponseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.contentLength += n

	return n, err
}

func (rw *statsResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *statsResponseWriter) Status() int {
	return rw.status
}

func (rw *statsResponseWriter) ContentLength() int {
	return rw.contentLength
}
