package mohttp

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
)

func Serve(c context.Context, handlers ...Handler) {
	beforeFuncs := []HandlerFunc{}
	mainFuncs := []HandlerFunc{}

	for _, h := range handlers {
		if b, ok := h.(PhaseHandler); ok {
			beforeFuncs = append(beforeFuncs, b.BeforeMain)
		}

		mainFuncs = append(mainFuncs, h.Handle)
	}

	var (
		idx   = 0
		funcs = append(beforeFuncs, mainFuncs...)
	)

	next := HandlerFunc(func(c context.Context) {
		if idx > len(funcs) {
			return
		}

		cur := funcs[idx]
		idx++

		cur(c)
	})

	c = nextStore.Set(c, next)
	next(c)
}

func HTTPContext(w http.ResponseWriter, r *http.Request, p httprouter.Params) context.Context {
	c := context.Background()
	c = WithRequest(c, r)
	c = WithResponseWriter(c, w)
	c = WithPathValues(c, &PathValues{
		Params: params(p),
		Query:  query(r.URL.Query()),
	})

	return c
}

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

func Next(c context.Context) {
	nextStore.Get(c).(HandlerFunc)(c)
}
