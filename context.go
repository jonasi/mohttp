package mohttp

import (
	"github.com/jonasi/mohttp/contextutil"
	"golang.org/x/net/context"
)

// ContextValueAccessors generates a pair of funcs to get and set values on the context. The first
// func returns a Handler that will set the value in a middleware chain.  The second func can be used
// to retrieve that value from a handler body.
func ContextValueAccessors(str string) (func(interface{}) PriorityHandler, func(context.Context) interface{}) {
	st := NewContextValueStore(str)

	return func(val interface{}) PriorityHandler {
		return PriorityHandlerFunc(-100, func(c context.Context) {
			c = st.Set(c, val)
			Next(c)
		})
	}, st.Get
}

// NewContextValueStore is a small wrapper around contextutil.ValueStore that turns the provided string
// into a unique key
func NewContextValueStore(key string) contextutil.ValueStore {
	return &contextValueStore{
		contextutil.NewValueStore(contextutil.NewKey(key)),
	}
}

type contextValueStore struct {
	contextutil.ValueStore
}
