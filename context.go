package mohttp

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
			c.Next.Handle(c)
		})
	}, st
}

// implementation pulled from https://go-review.googlesource.com/#/c/10910/
func NewContextKey(str string) interface{} {
	return (*contextKey)(&str)
}

func (k *contextKey) String() string { return string(*k) }
