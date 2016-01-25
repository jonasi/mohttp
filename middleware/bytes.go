package middleware

import (
	"bytes"
	"errors"
	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"sync"
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
		rd io.Reader
	)

	if b, ok := res.([]byte); ok {
		rd = bytes.NewReader(b)
	} else if b, ok := res.(BytesProvider); ok {
		rd = bytes.NewReader(b.Bytes())
	} else if r, ok := res.(io.Reader); ok {
		rd = r
	} else {
		return errors.New("DataProvider result cannot be coerced into []byte")
	}

	_, err := io.Copy(&onFirstWriteWriter{
		Writer: rw,
		OnFirstWrite: func(b []byte) {
			ct := http.DetectContentType(b)
			rw.Header().Set("Content-Type", ct)
		},
	}, rd)
	return err
}

type onFirstWriteWriter struct {
	io.Writer
	OnFirstWrite func([]byte)
	once         sync.Once
}

func (w *onFirstWriteWriter) Write(b []byte) (int, error) {
	w.once.Do(func() {
		w.OnFirstWrite(b)
	})

	return w.Writer.Write(b)
}
