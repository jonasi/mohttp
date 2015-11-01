package mohttp

import (
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"time"
)

type RequestStats struct {
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

func Logger(fn func(*RequestStats)) Handler {
	return HandlerFunc(func(c context.Context) {
		req := GetRequest(c)

		st := &RequestStats{
			StartTime: time.Now(),
			Protocol:  req.Proto,
			Method:    req.Method,
			URL:       req.URL,
			UserAgent: req.UserAgent(),
			ClientIP:  req.RemoteAddr,
			Referer:   req.Referer(),
		}

		rw := NewResponseWriter(GetResponseWriter(c))
		c = WithResponseWriter(c, rw)

		GetNext(c).Handle(c)

		st.Duration = time.Since(st.StartTime)
		st.StatusCode = rw.Status()
		st.ContentLength = rw.ContentLength()

		fn(st)
	})
}

func NewResponseWriter(rw http.ResponseWriter) *StatsResponseWriter {
	return &StatsResponseWriter{
		ResponseWriter: rw,
		status:         200,
	}
}

type StatsResponseWriter struct {
	http.ResponseWriter
	status        int
	contentLength int
}

func (rw *StatsResponseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.contentLength += n

	return n, err
}

func (rw *StatsResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *StatsResponseWriter) Status() int {
	return rw.status
}

func (rw *StatsResponseWriter) ContentLength() int {
	return rw.contentLength
}
