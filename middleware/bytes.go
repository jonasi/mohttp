package middleware

import (
	"errors"
	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
	"net/http"
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

	mohttp.GetResponseWriter(c).WriteHeader(500)
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

type DetectTypeResponder struct {
	OnError func(context.Context, error)
}

func (r *DetectTypeResponder) HandleErr(c context.Context, err error) {
	if r.OnError != nil {
		r.OnError(c, err)
		return
	}

	mohttp.GetResponseWriter(c).WriteHeader(500)
}

func (r *DetectTypeResponder) HandleResult(c context.Context, res interface{}) error {
	var (
		rw = mohttp.GetResponseWriter(c)
		b  []byte
	)

	if bytes, ok := res.([]byte); ok {
		b = bytes
	} else if bytes, ok := res.(BytesProvider); ok {
		b = bytes.Bytes()
	} else {
		return errors.New("DataProvider result cannot be coerced into []byte")
	}

	ct := http.DetectContentType(b)
	rw.Header().Set("Content-Type", ct)

	rw.Write(b)
	return nil
}
