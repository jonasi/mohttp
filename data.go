package mohttp

import (
	"golang.org/x/net/context"
)

var dataHandlerContext = NewContextValueStore("github.com/jonasi/http.DataHandler")

type DataResponder interface {
	HandleErr(context.Context, error)
	HandleResult(context.Context, interface{}) error
}

func WithResponder(c context.Context, r DataResponder) context.Context {
	return dataHandlerContext.Set(c, r)
}

type DataHandlerFunc func(c context.Context) (interface{}, error)

func DataHandler(fn DataHandlerFunc) Handler {
	return HandlerFunc(func(c context.Context) {
		r, ok := dataHandlerContext.Get(c).(DataResponder)

		if !ok {
			panic("No data responder set")
		}

		var (
			err  error
			data interface{}
		)

		defer func() {
			if err2 := recoverErr(); err2 != nil {
				err = err2
			}

			if err == nil {
				err = r.HandleResult(c, data)
			}

			if err != nil {
				r.HandleErr(c, err)
			}
		}()

		data, err = fn(c)
	})
}
