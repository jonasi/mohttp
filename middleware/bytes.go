package middleware

import (
	"errors"
	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
)

type BytesProvider interface {
	Bytes() []byte
}

type BytesResponder struct {
	OnError func(context.Context, error)
}

func (r *BytesResponder) HandleErr(c context.Context, err error) {
	if r.OnError != nil {
		r.OnError(c, err)
		return
	}

	mohttp.Error(c, err.Error(), 500)
}

func (r *BytesResponder) HandleResult(c context.Context, res interface{}) error {
	rw := mohttp.GetResponseWriter(c)
	rw.Header().Set("Content-Type", "application/octet-stream")

	if b, ok := res.([]byte); ok {
		rw.Write(b)
		return nil
	}

	if b, ok := res.(BytesProvider); ok {
		rw.Write(b.Bytes())
		return nil
	}

	return errors.New("DataProvider result cannot be coerced into []byte")
}
