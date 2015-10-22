package http

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Params  httprouter.Params
	Context context.Context
	Next    Handler
}

func (c *Context) QueryString(k string) string {
	return c.Request.URL.Query().Get(k)
}

func (c *Context) QueryInt(k string) int {
	v := c.QueryString(k)
	iv, _ := strconv.Atoi(v)

	return iv
}
