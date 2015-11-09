package middleware

import (
	"errors"

	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
)

type Stringer interface {
	String() string
}

type TextResponder struct {
	OnError func(context.Context, error)
}

func (t *TextResponder) HandleErr(c context.Context, err error) {
	if t.OnError != nil {
		t.OnError(c, err)
		return
	}

	mohttp.GetResponseWriter(c).WriteHeader(500)
}

func (t *TextResponder) HandleResult(c context.Context, res interface{}) error {
	rw := mohttp.GetResponseWriter(c)
	rw.Header().Set("Content-Type", "text/plain")

	if str, ok := res.(string); ok {
		rw.Write([]byte(str))
		return nil
	}

	if str, ok := res.(Stringer); ok {
		rw.Write([]byte(str.String()))
		return nil
	}

	return errors.New("DataProvider result cannot be coerced into string")
}

func TextHandler(h mohttp.DataHandlerFunc) mohttp.Handler {
	return mohttp.HandlerFunc(func(c context.Context) {
		r, _ := mohttp.GetDataResponder(c)

		if _, ok := r.(*TextResponder); !ok {
			c = mohttp.WithDataResponder(c, &TextResponder{})
		}

		h.Handle(c)
	})
}
