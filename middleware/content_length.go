package middleware

import (
	"bytes"
	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"sync"
)

type ContentLengthHandler struct {
}

func (h *ContentLengthHandler) Handle(c context.Context) {
	old := mohttp.GetResponseWriter(c)
	rw := clrw{ResponseWriter: old}

	c = mohttp.WithResponseWriter(c, &rw)
	mohttp.Next(c)

	if rw.useBuf {
		rw.Header().Set("Content-Length", strconv.Itoa(rw.buf.Len()))
		rw.buf.WriteTo(old)
	}
}

type clrw struct {
	http.ResponseWriter
	buf    bytes.Buffer
	once   sync.Once
	useBuf bool
}

func (c *clrw) check() {
	if c.Header().Get("Content-Length") != "" {
		c.useBuf = false
	} else if c.Header().Get("Transfer-Encoding") != "" {
		c.useBuf = false
	} else {
		c.useBuf = true
	}
}

func (c *clrw) Write(b []byte) (int, error) {
	c.once.Do(c.check)

	if c.useBuf {
		return c.buf.Write(b)
	}

	return c.ResponseWriter.Write(b)
}
