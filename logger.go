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
		st := &RequestStats{
			StartTime: time.Now(),
			Protocol:  c.Request.Proto,
			Method:    c.Request.Method,
			URL:       c.Request.URL,
			UserAgent: c.Request.UserAgent(),
			ClientIP:  c.Request.RemoteAddr,
			Referer:   c.Request.Referer(),
		}

		rw := NewResponseWriter(c.Writer)
		c.Writer = rw

		c.Next.Handle(c)

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
