package mohttp

import (
	"golang.org/x/net/context"
)

var dataHandlerContext = NewContextValueStore("github.com/jonasi/http.DataHandler")

type DataResponder interface {
	HandleErr(context.Context, error)
	HandleResult(context.Context, interface{}) error
}

func GetDataResponder(c context.Context) (DataResponder, bool) {
	r, ok := dataHandlerContext.Get(c).(DataResponder)
	return r, ok
}

func WithDataResponder(c context.Context, r DataResponder) context.Context {
	return dataHandlerContext.Set(c, r)
}

func DataResponderHandler(d DataResponder) Handler {
	return HandlerFunc(func(c context.Context) {
		c = WithDataResponder(c, d)
		Next(c)
	})
}

type DataHandlerFunc func(context.Context) (interface{}, error)

func (fn DataHandlerFunc) Handle(c context.Context) {
	r, ok := dataHandlerContext.Get(c).(DataResponder)

	if !ok {
		panic("No DataResponder set in handler chain")
	}

	var (
		err    error
		result interface{}
	)

	defer func() {
		if err2 := recoverErr(); err2 != nil {
			err = err2
		}

		if err == nil {
			err = r.HandleResult(c, result)
		}

		if err != nil {
			r.HandleErr(c, err)
		}
	}()

	result, err = fn(c)
}
