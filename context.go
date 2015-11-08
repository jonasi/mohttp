package mohttp

import (
	"github.com/jonasi/mohttp/contextutil"
	"golang.org/x/net/context"
)

func ContextValueAccessors(str string) (func(interface{}) PriorityHandler, func(context.Context) interface{}) {
	st := NewContextValueStore(str)

	return func(val interface{}) PriorityHandler {
		return PriorityHandlerFunc(-100, func(c context.Context) {
			c = st.Set(c, val)
			Next(c)
		})
	}, st.Get
}

func NewContextValueStore(key string) contextutil.ValueStore {
	return &contextValueStore{
		contextutil.NewValueStore(contextutil.NewKey(key)),
	}
}

type contextValueStore struct {
	contextutil.ValueStore
}
