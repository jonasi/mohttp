package mohttp

import (
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
	return HandlerFunc(func(c *Context) {
		req := c.Request()

		st := &RequestStats{
			StartTime: time.Now(),
			Protocol:  req.Proto,
			Method:    req.Method,
			URL:       req.URL,
			UserAgent: req.UserAgent(),
			ClientIP:  req.RemoteAddr,
			Referer:   req.Referer(),
		}

		rw := NewResponseWriter(c.ResponseWriter())
		c = c.WithResponseWriter(rw)

		c.Next().Handle(c)

		st.Duration = time.Since(st.StartTime)
		st.StatusCode = rw.Status()
		st.ContentLength = rw.ContentLength()

		fn(st)
	})
}

func NewResponseWriter(rw http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: rw,
		status:         200,
	}
}

type ResponseWriter struct {
	http.ResponseWriter
	status        int
	contentLength int
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.contentLength += n

	return n, err
}

func (rw *ResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *ResponseWriter) Status() int {
	return rw.status
}

func (rw *ResponseWriter) ContentLength() int {
	return rw.contentLength
}
