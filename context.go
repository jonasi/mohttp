package mohttp

import (
	"github.com/jonasi/mohttp/contextutil"
	"golang.org/x/net/context"
	"net/http"
)

var (
	reqStore  = NewContextValueStore("github.com/jonasi/mohttp.Request")
	respStore = NewContextValueStore("github.com/jonasi/mohttp.ResponseWriter")
	pathStore = NewContextValueStore("github.com/jonasi/mohttp.PathValues")
	nextStore = NewContextValueStore("github.com/jonasi/mohttp.Next")
)

func WithRequest(c context.Context, r *http.Request) context.Context {
	return reqStore.Set(c, r)
}

func GetRequest(c context.Context) *http.Request {
	return reqStore.Get(c).(*http.Request)
}

func WithResponseWriter(c context.Context, w http.ResponseWriter) context.Context {
	return respStore.Set(c, w)
}

func GetResponseWriter(c context.Context) http.ResponseWriter {
	return respStore.Get(c).(http.ResponseWriter)
}

func WithPathValues(c context.Context, p *PathValues) context.Context {
	return pathStore.Set(c, p)
}

func GetPathValues(c context.Context) *PathValues {
	return pathStore.Get(c).(*PathValues)
}

func WithNext(c context.Context, h Handler) context.Context {
	return nextStore.Set(c, h)
}

func GetNext(c context.Context) Handler {
	return nextStore.Get(c).(Handler)
}

func NewContextValuePair(str string) (func(interface{}) Handler, contextutil.ValueStore) {
	st := NewContextValueStore(str)

	return func(val interface{}) Handler {
		return HandlerFunc(func(c context.Context) {
			c = st.Set(c, val)
			GetNext(c).Handle(c)
		})
	}, st
}

func NewContextValueStore(key string) contextutil.ValueStore {
	return &contextValueStore{
		contextutil.NewValueStore(contextutil.NewKey(key)),
	}
}

type contextValueStore struct {
	contextutil.ValueStore
}
