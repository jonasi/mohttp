package mohttp

import (
	"golang.org/x/net/context"
)

var dataHandlerContext = NewContextValueStore("github.com/jonasi/http.DataHandler")

// A DataResponder is responsible for formatting the return values of a DataHandlerFunc
// onto the http Response
type DataResponder interface {
	HandleErr(context.Context, error)
	HandleResult(context.Context, interface{}) error
}

// GetDataResponder retrieves the current DataResponder, if it exists, from the context
func GetDataResponder(c context.Context) (DataResponder, bool) {
	r, ok := dataHandlerContext.Get(c).(DataResponder)
	return r, ok
}

// WithDataResponder returns a new context with the provided DataResponder set
func WithDataResponder(c context.Context, r DataResponder) context.Context {
	return dataHandlerContext.Set(c, r)
}

// DataResponderHandler returns a Handler whose single responsibility is to set
// the provided DataResponder on the context
func DataResponderHandler(d DataResponder) Handler {
	return HandlerFunc(func(c context.Context) {
		c = WithDataResponder(c, d)
		Next(c)
	})
}

// A DataHandlerFunc is a Handle which separates out the logic of
// retrieving/returning data from the logic of formatting/serializing it
// on the wire
type DataHandlerFunc func(context.Context) (interface{}, error)

// Handle satisfies the Handler interface.  If no DataResponder has been
// set on the context, this method will panic.
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

func StaticDataHandler(v interface{}) DataHandlerFunc {
	return DataHandlerFunc(func(c context.Context) (interface{}, error) {
		return v, nil
	})
}
