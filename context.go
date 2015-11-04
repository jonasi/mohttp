package mohttp

import (
	"github.com/jonasi/mohttp/contextutil"
	"golang.org/x/net/context"
)

func NewContextValuePair(str string) (func(interface{}) PhaseHandler, contextutil.ValueStore) {
	st := NewContextValueStore(str)

	return func(val interface{}) PhaseHandler {
		return BeforeHandlerFunc(func(c context.Context) {
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
