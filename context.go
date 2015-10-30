package mohttp

import (
	"golang.org/x/net/context"
	"net/http"
)

var (
	reqKey  = NewContextKey("github.com/jonasi/mohttp.Request")
	respKey = NewContextKey("github.com/jonasi/mohttp.ResponseWriter")
	pathKey = NewContextKey("github.com/jonasi/mohttp.PathValues")
	nextKey = NewContextKey("github.com/jonasi/mohttp.Next")
)

type Context struct {
	context.Context
}

func (c *Context) WithRequest(r *http.Request) *Context {
	return &Context{context.WithValue(c, reqKey, r)}
}

func (c *Context) Request() *http.Request {
	return c.Value(reqKey).(*http.Request)
}

func (c *Context) WithResponseWriter(w http.ResponseWriter) *Context {
	return &Context{context.WithValue(c, respKey, w)}
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.Value(respKey).(http.ResponseWriter)
}

func (c *Context) WithPathValues(p *PathValues) *Context {
	return &Context{context.WithValue(c, pathKey, p)}
}

func (c *Context) PathValues() *PathValues {
	return c.Value(pathKey).(*PathValues)
}

func (c *Context) WithNext(h Handler) *Context {
	return &Context{context.WithValue(c, nextKey, h)}
}

func (c *Context) Next() Handler {
	return c.Value(nextKey).(Handler)
}

type contextKey string

type ContextValueStore interface {
	Get(*Context) interface{}
	Set(*Context, interface{}) *Context
}

func NewContextValueStore(str string) ContextValueStore {
	return &contextValueStore{contextKey(str)}
}

type contextValueStore struct {
	k contextKey
}

func (c *contextValueStore) Get(ctxt *Context) interface{} {
	return ctxt.Context.Value(c.k)
}

func (c *contextValueStore) Set(ctxt *Context, val interface{}) *Context {
	ctxt.Context = context.WithValue(ctxt.Context, c.k, val)
	return ctxt
}

func NewContextValueMiddleware(str string) (func(interface{}) Handler, ContextValueStore) {
	st := NewContextValueStore(str)

	return func(val interface{}) Handler {
		return HandlerFunc(func(c *Context) {
			c = st.Set(c, val)
			c.Next().Handle(c)
		})
	}, st
}

// implementation pulled from https://go-review.googlesource.com/#/c/10910/
func NewContextKey(str string) interface{} {
	return (*contextKey)(&str)
}

func (k *contextKey) String() string { return string(*k) }
