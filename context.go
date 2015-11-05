package mohttp

import (
	"github.com/jonasi/mohttp/contextutil"
	"golang.org/x/net/context"
)

func NewContextValuePair(str string) (func(interface{}) PriorityHandler, contextutil.ValueStore) {
	st := NewContextValueStore(str)

	return func(val interface{}) PriorityHandler {
		return PriorityHandlerFunc(-100, func(c context.Context) {
			c = st.Set(c, val)
			Next(c)
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
