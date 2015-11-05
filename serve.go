package mohttp

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
	"sort"
)

func prio(h Handler) int {
	if ph, ok := h.(PriorityHandler); ok {
		return ph.Priority()
	}

	return 0
}

type sortedHandlers []Handler

func (s sortedHandlers) Len() int           { return len(s) }
func (s sortedHandlers) Swap(a, b int)      { s[a], s[b] = s[b], s[a] }
func (s sortedHandlers) Less(a, b int) bool { return prio(s[a]) < prio(s[b]) }

func Serve(c context.Context, handlers ...Handler) {
	var (
		idx    = 0
		sorted = sortedHandlers(handlers)
	)

	sort.Stable(sorted)

	next := HandlerFunc(func(c context.Context) {
		if idx >= len(sorted) {
			return
		}

		cur := sorted[idx]
		idx++

		cur.Handle(c)
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
